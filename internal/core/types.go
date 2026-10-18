package core

// ===== DTO / API-контракты Nubes API =====
//
// Здесь только структуры (без логики):
//   - generic*  — простые payload'ы создания инстанса/операции/параметра
//   - universal* — ответы схемы операций (cfsParams) и параметры
//   - InstanceStateResponse / ApiOperation — состояние инстанса
//   - operationStatusResponse / opStage — статус операции при поллинге

type genericInstanceReq struct {
	ServiceId   int    `json:"serviceId"`
	DisplayName string `json:"displayName"`
	Descr       string `json:"descr"`
}

type genericOpReq struct {
	InstanceUid string `json:"instanceUid"`
	Operation   string `json:"operation"`
}

type genericParamReq struct {
	InstanceOperationUid   string `json:"instanceOperationUid"`
	SvcOperationCfsParamId int    `json:"svcOperationCfsParamId"`
	ParamValue             string `json:"paramValue"`
}

// ===== UNIVERSAL FLOW V6 (baseline) =====

type universalOpResponse struct {
	InstanceOperation universalOperation `json:"instanceOperation"`
}

type universalOperation struct {
	CfsParams []universalCfsParam `json:"cfsParams"`
}

// universalOpDefaultResponse — схема операции из GET /instanceOperations/default/{opId}.
// Не вычисляет instance-специфичные выражения (getResourceRealmConfig) и не падает
// с 500 на проблемных инстансах (см. NOTES/30_analysis/DEBUG_REPORT_VC_VDC_500.md).
type universalOpDefaultResponse struct {
	SvcOperation universalOperation `json:"svcOperation"`
}

type universalSubParam struct {
	DefaultValue            string `json:"defaultValue"`
	DataType                string `json:"dataType"`
	SvcOperationCfsSubparam string `json:"svcOperationCfsSubparam"`
	IsRequired              bool   `json:"isRequired"`
}

type universalCfsParam struct {
	SvcOperationCfsParamId int                          `json:"svcOperationCfsParamId"`
	ParamValue             *string                      `json:"paramValue"`
	DefaultValue           *string                      `json:"defaultValue"`
	IsRequired             bool                         `json:"isRequired"`
	DataType               string                       `json:"dataType"`
	Name                   string                       `json:"name"`
	Code                   string                       `json:"code"`
	Label                  string                       `json:"label"`
	SvcOperationCfsParam   string                       `json:"svcOperationCfsParam"`
	RefSvcId               *int                         `json:"refSvcId"`
	DataDescriptor         map[string]universalSubParam `json:"dataDescriptor"`
}

// ===== Instance state =====

type ApiOperation struct {
	SvcOperationId int    `json:"svcOperationId"`
	Operation      string `json:"operation"`
}

type InstanceStateResponse struct {
	InstanceUid     string `json:"instanceUid"`
	ServiceId       int    `json:"serviceId"`
	ExplainedStatus string `json:"explainedStatus"`
	IsDeleted       bool   `json:"isDeleted"`
	// IsCreated — указатель: false означает "инстанс не создан" (полный откат операции),
	// а nil — что поле не пришло (legacy-эндпоинт). Различать обязательно, иначе
	// отсутствующее поле будет ложно трактоваться как "не создан".
	IsCreated             *bool          `json:"isCreated"`
	OperationIsInProgress bool           `json:"operationIsInProgress"`
	OperationIsPending    bool           `json:"operationIsPending"`
	AvailableOperations   []ApiOperation `json:"availableOperations"`
}

// ===== Operation status (polling) =====

type opStage struct {
	InstanceOperationStageUid string  `json:"instanceOperationStageUid"`
	Stage                     string  `json:"stage"`
	IsSuccessful              bool    `json:"isSuccessful"`
	DtFinish                  *string `json:"dtFinish"`
	Duration                  float64 `json:"duration"`
	StageMsg                  *string `json:"stageMsg"`
}

type operationStatusResponse struct {
	InstanceOperation struct {
		DtFinish     *string   `json:"dtFinish"`
		IsSuccessful *bool     `json:"isSuccessful"`
		ErrorLog     *string   `json:"errorLog"`
		IsInProgress bool      `json:"isInProgress"`
		IsPending    bool      `json:"isPending"`
		Duration     *float64  `json:"duration"`
		Stages       []opStage `json:"stages"`
	} `json:"instanceOperation"`
}
