package main

var ParamMaxes = []float64{9999.0, 2000.0, 255.0, 255.0, 255.0, 255.0, 500.0, 500.0}

// Class はツクールMVにおける職業
type Class struct {
	ExpParams []float64 `json:"expParams"`
	ID        float64   `json:"id"`
	Learnings []struct {
		Level   float64 `json:"level"`
		Note    string  `json:"note"`
		SkillID float64 `json:"skillId"`
	} `json:"learnings"`
	Name   string      `json:"name"`
	Note   string      `json:"note"`
	Params [][]float64 `json:"params"`
	Traits []struct {
		Code   float64 `json:"code"`
		DataID float64 `json:"dataId"`
		Value  float64 `json:"value"`
	} `json:"traits"`
}

// FetchLastParams はParamsの最後の値を取り出したスライスを返す
// Classes.jsonはの一番最後の値でチャートを生成するため。
func (c Class) FetchLastParams() (f []float64) {
	f = make([]float64, len(c.Params))
	for i, v := range c.Params {
		p := v[len(v)-1]
		f[i] = p
	}
	return
}

// Classes はClasses.json構造体
type Classes []Class
