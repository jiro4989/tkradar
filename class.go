package main

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

// Classes はClasses.json構造体
type Classes []Class
