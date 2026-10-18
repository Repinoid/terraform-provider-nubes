package resources_core

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/repinoid/terraform-provider-nubes/internal/core"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ExtractUserPassword вытаскивает пароль пользователя из JSON-значения секрета Vault
// родительского инстанса.
//
// ЗАЧЕМ: у части сервисов (проверено на postgres) платформа САМА генерирует пароль
// в операции create_user и кладёт его в секрет Vault родителя. Формат значения
// (проверено через GET /instances/<uid>/vault/users):
//
//	{"<username>": {"password": "<пароль>"}}
//
// Подресурс обязан вернуть этот пароль СВОИМ выходом сразу после create_user, иначе
// потребитель не сможет получить его в том же apply (vault_secrets родителя — Computed
// и обновляется только при Read родителя). См. TOOLS/ARCHITECTURE.md,
// раздел «Subresource-born secrets».
//
// Возвращает "" без ошибки, если пользователя в секрете нет или пароль отсутствует:
// это не фатально (пароль остаётся неизвестным до следующего чтения) и не должно
// ломать create.
func ExtractUserPassword(secretJSON string, username string) (string, error) {
	raw := strings.TrimSpace(secretJSON)
	if raw == "" {
		return "", nil
	}
	name := strings.TrimSpace(username)
	if name == "" {
		return "", fmt.Errorf("пустое имя пользователя для поиска пароля в Vault")
	}

	// Основной формат: {"<username>": {"password": "..."}}.
	var byUser map[string]struct {
		Password string `json:"password"`
	}
	if err := json.Unmarshal([]byte(raw), &byUser); err == nil {
		if entry, ok := byUser[name]; ok && strings.TrimSpace(entry.Password) != "" {
			return strings.TrimSpace(entry.Password), nil
		}
		// Регистр имени может отличаться (платформа нормализует его сама).
		for key, entry := range byUser {
			if strings.EqualFold(strings.TrimSpace(key), name) && strings.TrimSpace(entry.Password) != "" {
				return strings.TrimSpace(entry.Password), nil
			}
		}
	}

	// Упрощённый формат: {"password": "..."} — единственный пользователь в секрете.
	var flat struct {
		Password string `json:"password"`
	}
	if err := json.Unmarshal([]byte(raw), &flat); err == nil {
		if strings.TrimSpace(flat.Password) != "" {
			return strings.TrimSpace(flat.Password), nil
		}
	}

	return "", nil
}

// ResolveUserPasswordFromVault читает пароль пользователя из секретов Vault
// родительского инстанса и возвращает готовое значение для атрибута password.
//
// Зачем функция, а не код в шаблоне: подресурс сохраняет состояние в НЕСКОЛЬКИХ
// местах (обычный create, усыновление уже существующего объекта, подтверждение по
// state_out). Если хотя бы в одной ветке оставить password неизвестным, Terraform
// падает с «Provider returned invalid result object after apply: ... still indicated
// an unknown value». Поэтому все ветки обязаны проставлять КОНКРЕТНОЕ значение —
// через эту функцию.
//
// Поведение:
//   - удалось прочитать пароль → конкретное значение;
//   - не удалось → возвращается current, если оно уже известно (например, из state),
//     иначе типизированный null (null тоже «known» для Terraform).
//
// Второй результат — текст предупреждения: пустая строка, если всё хорошо либо у
// родителя вообще нет секретов в Vault (для части сервисов это норма). Почему не
// просто тихо вернуть null: иначе "пароль не прочитан" выглядел бы в apply как
// успешный пустой пароль, и причина (недоступен API/Vault, нет имени, нет записи)
// осталась бы невидимой.
func ResolveUserPasswordFromVault(ctx context.Context, client *core.UniversalClient, instanceUID string, username string, current types.String) (types.String, string) {
	fallback := func() types.String {
		if current.IsNull() || current.IsUnknown() {
			return types.StringNull()
		}
		return current
	}

	if client == nil || strings.TrimSpace(instanceUID) == "" {
		return fallback(), ""
	}

	details, err := client.GetInstanceStateDetails(ctx, instanceUID)
	if err != nil {
		return fallback(), fmt.Sprintf("не удалось прочитать состояние родительского инстанса %s: %v", instanceUID, err)
	}
	if details == nil || len(details.Vault.Fields) == 0 {
		// У родителя нет секретов Vault — не ошибка.
		return fallback(), ""
	}

	secrets, serr := client.GetInstanceVaultSecrets(ctx, instanceUID, details.Vault.Fields)
	if serr != nil {
		return fallback(), fmt.Sprintf("не удалось прочитать секреты Vault %v у инстанса %s: %v", details.Vault.Fields, instanceUID, serr)
	}

	for _, secretValue := range secrets {
		pw, perr := ExtractUserPassword(secretValue, username)
		if perr != nil {
			continue
		}
		if pw != "" {
			return types.StringValue(pw), ""
		}
	}

	return fallback(), fmt.Sprintf("пароль пользователя %q не найден в секретах Vault %v (инстанс %s)", username, details.Vault.Fields, instanceUID)
}
