package analyzer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/diploma/worker-static-analyzer/internal/model"
)

type Analyzer struct {
	binaryPath string
}

func New(binaryPath string) *Analyzer {
	return &Analyzer{binaryPath: binaryPath}
}

type analyzerConfig struct {
	Input        string      `json:"input"`
	Output       string      `json:"output"`
	OutputFormat string      `json:"output_format"`
	Analysis     analysisCfg `json:"analysis"`
	Debug        debugCfg    `json:"debug"`
	Features     featuresCfg `json:"features"`
}

type analysisCfg struct {
	MaxLoopDepth        int  `json:"max_loop_depth"`
	AnalyzeDependencies bool `json:"analyze_dependencies"`
	AnalyzeSCEV         bool `json:"analyze_scev"`
}

type debugCfg struct {
	Verbose    bool `json:"verbose"`
	DumpLoops  bool `json:"dump_loops"`
	DumpSCEV   bool `json:"dump_scev"`
	DumpMemory bool `json:"dump_memory"`
}

type featuresCfg struct {
	EnableFingerprint    bool `json:"enable_fingerprint"`
	EnableClassification bool `json:"enable_classification"`
}

func (a *Analyzer) Run(ctx context.Context, sourceFile, workDir string) ([]model.Pattern, error) {
	confPath := filepath.Join(workDir, "conf.json")
	outPath := filepath.Join(workDir, "out.json")

	cfg := analyzerConfig{
		Input:        filepath.Base(sourceFile),
		Output:       "out.json",
		OutputFormat: "json",
		Analysis: analysisCfg{
			MaxLoopDepth:        4,
			AnalyzeDependencies: true,
			AnalyzeSCEV:         true,
		},
		Debug:    debugCfg{},
		Features: featuresCfg{EnableFingerprint: true, EnableClassification: true},
	}

	confBytes, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal conf: %w", err)
	}
	if err := os.WriteFile(confPath, confBytes, 0o644); err != nil {
		return nil, fmt.Errorf("write conf.json: %w", err)
	}

	cmd := exec.CommandContext(ctx, a.binaryPath, "conf.json", "--quiet")
	cmd.Dir = workDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("run analyzer (%s): %w", a.binaryPath, err)
	}

	raw, err := os.ReadFile(outPath)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", outPath, err)
	}

	var patterns []model.Pattern
	if err := json.Unmarshal(raw, &patterns); err != nil {
		return nil, fmt.Errorf("parse out.json: %w", err)
	}
	return patterns, nil
}
