package model

type DataPath struct {
	TarPath  string `json:"TarPath"`
	NumFile  int    `json:"NumFile"`
	SumBytes int    `json:"SumBytes"`
	AvgBytes int    `json:"AvgBytes"`
}

type DataPathList []DataPath

type DatadataPost struct {
	DataFileLimit  DataPathList `json:"DataFileLimit"`
	DataBytesLimit DataPathList `json:"DataBytesLimit"`
}
