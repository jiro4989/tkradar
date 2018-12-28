package main

// Class はツクールMVにおける職業
type Class struct {
	ExpParams []float32 `json:"expParams"`
	ID        float32   `json:"id"`
	Learnings []struct {
		Level   float32 `json:"level"`
		Note    string  `json:"note"`
		SkillID float32 `json:"skillId"`
	} `json:"learnings"`
	Name   string      `json:"name"`
	Note   string      `json:"note"`
	Params [][]float32 `json:"params"`
	Traits []struct {
		Code   float32 `json:"code"`
		DataID float32 `json:"dataId"`
		Value  float32 `json:"value"`
	} `json:"traits"`
}

// Classes はClasses.json構造体
type Classes []Class
