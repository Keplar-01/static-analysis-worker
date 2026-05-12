package model

// Pattern — один паттерн доступа к памяти, выданный LLVM-анализатором (C++).
type Pattern struct {
	SequenceIndex      int      `json:"sequence_index"`
	AccessKind         string   `json:"access_kind"`
	Affine             bool     `json:"affine"`
	Alignment          *int     `json:"alignment"`
	BaseKind           string   `json:"base_kind"`
	BaseSymbol         string   `json:"base_symbol"`
	Conditional        bool     `json:"conditional"`
	ContiguousBlock    *int     `json:"contiguous_block"`
	Dependence         string   `json:"dependence"`
	Depth              int      `json:"depth"`
	FillFactor         float64  `json:"fill_factor"`
	Function           string   `json:"function"`
	HasIndexedAddr     bool     `json:"has_indexed_addressing"`
	IndexedByMemory    bool     `json:"indexed_by_memory"`
	LoadCount          int      `json:"load_count"`
	PatternFingerprint string   `json:"pattern_fingerprint"`
	PatternSig         string   `json:"pattern_signature"`
	PatternType        string   `json:"pattern_type"`
	SourceColumn       int      `json:"source_column"`
	SourceFile         string   `json:"source_file"`
	SourceLine         int      `json:"source_line"`
	StoreCount         int      `json:"store_count"`
	Stride             *float64 `json:"stride"`
	WorkingSetBytes    int      `json:"working_set_bytes"`
}
