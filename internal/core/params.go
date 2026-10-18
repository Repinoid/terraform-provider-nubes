package core

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ===== Нормализация значений CFS-параметров =====

// buildMapFixedDefault строит JSON-объект из дефолтов sub-параметров map-fixed.
func buildMapFixedDefault(param universalCfsParam) string {
	result := make(map[string]string, len(param.DataDescriptor))
	for key, sub := range param.DataDescriptor {
		result[key] = sub.DefaultValue
	}
	b, err := json.Marshal(result)
	if err != nil {
		return "{}"
	}
	return string(b)
}

func normalizeUniversalValueV6(val string, param universalCfsParam) string {
	trimmed := strings.TrimSpace(val)
	if strings.EqualFold(trimmed, "null") {
		trimmed = ""
	}
	if trimmed == "\"\"" {
		trimmed = ""
	}

	// map-fixed с DataDescriptor: если значение пустое или "{}" — строим JSON из дефолтов sub-params.
	dataType := strings.ToLower(param.DataType)
	if (dataType == "map-fixed" || strings.HasPrefix(dataType, "map")) && len(param.DataDescriptor) > 0 {
		if trimmed == "" || trimmed == "{}" {
			return buildMapFixedDefault(param)
		}
		return trimmed
	}

	if trimmed != "" {
		return trimmed
	}

	nameHint := strings.ToLower(param.Name + " " + param.Code + " " + param.Label + " " + param.SvcOperationCfsParam)

	if strings.Contains(dataType, "array") || strings.Contains(nameHint, "array") || strings.Contains(nameHint, "list") {
		return "[]"
	}
	if strings.Contains(dataType, "map") || strings.Contains(dataType, "json") || strings.Contains(nameHint, "map") || strings.Contains(nameHint, "json") {
		return "{}"
	}
	if strings.Contains(dataType, "integer") || strings.Contains(dataType, "int") {
		return "0"
	}
	if strings.Contains(dataType, "boolean") || strings.Contains(dataType, "bool") {
		return "false"
	}

	return trimmed
}

// paramLabelByID отдаёт человекочитаемое имя параметра операции по его ID.
//
// Зачем: в текстах ошибок числовой ID («не удалось установить параметр 779»)
// юзеру бесполезен — он не знает, что 779 это recordType. С именем из спеки
// ошибка становится самодостаточной.
func paramLabelByID(params []universalCfsParam, paramID int) string {
	for _, p := range params {
		if p.SvcOperationCfsParamId == paramID {
			return paramLabel(p)
		}
	}
	return fmt.Sprintf("id %d", paramID)
}

// paramLabel собирает имя параметра: Code (имя в API/спеке) + Name/Label (UI).
// ⚠️ ЗНАЧЕНИЕ параметра сюда НЕ попадает намеренно — в параметрах бывают секреты.
func paramLabel(p universalCfsParam) string {
	code := strings.TrimSpace(p.Code)
	name := strings.TrimSpace(p.Name)
	if name == "" {
		name = strings.TrimSpace(p.Label)
	}
	switch {
	case code != "" && name != "" && !strings.EqualFold(code, name):
		return fmt.Sprintf("%s (%s)", code, name)
	case code != "":
		return code
	case name != "":
		return name
	default:
		return fmt.Sprintf("id %d", p.SvcOperationCfsParamId)
	}
}
