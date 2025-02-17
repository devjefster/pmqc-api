package models

type PMQCData struct {
	UF map[string]State `json:"UF"`
}

type State struct {
	Municipios map[string]Municipio `json:"municipios"`
}

type Municipio struct {
	Amostras map[string]Amostra `json:"amostras"`
}

type Amostra struct {
	IdNumeric    int               `json:"IdNumeric"`
	DataColeta   string            `json:"DataColeta"`
	GrupoProduto string            `json:"GrupoProduto"`
	Produto      string            `json:"Produto"`
	Posto        Posto             `json:"Posto"`
	Ensaios      map[string]Ensaio `json:"Ensaios"`
}

type Posto struct {
	RazaoSocial   string  `json:"razaosocial"`
	CNPJ          string  `json:"CNPJ"`
	Distribuidora string  `json:"Distribuidora"`
	Endereco      string  `json:"endereco"`
	Complemento   string  `json:"complemento"`
	Bairro        string  `json:"bairro"`
	Latitude      float64 `json:"Latitude,string"`
	Longitude     float64 `json:"Longitude,string"`
}

type Ensaio struct {
	Resultado string `json:"Resultado"`
	Unidade   string `json:"Unidade"`
	Conforme  string `json:"Conforme"`
}
