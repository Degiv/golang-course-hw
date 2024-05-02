package main

type Opt int

const (
	size Opt = iota
	char
	color
)

type OptValue interface{}

type OptsMap map[Opt]OptValue

func Builder() *OptsMap {
	return &OptsMap{size: 15, char: 'X', color: 0} //default values
}

func (optsMap *OptsMap) SetSize(v int) *OptsMap {
	(*optsMap)[size] = v
	return optsMap
}

func (optsMap *OptsMap) SetChar(v rune) *OptsMap {
	(*optsMap)[char] = v
	return optsMap
}

func (optsMap *OptsMap) SetColor(v int) *OptsMap {
	(*optsMap)[color] = v
	return optsMap
}

func (optsMapPointer *OptsMap) Draw() {
	optsMap := *optsMapPointer
	println(optsMap[size].(int), string(optsMap[char].(rune)), optsMap[color].(int))
}

func main() {
	Builder().SetSize(3).SetChar('Y').SetColor(0).Draw()
}
