// cmd/seed/main.go — populates the database with realistic sample data.
//
// Usage (from the backend/ directory):
//
//	go run ./cmd/seed/main.go
//
// Environment variables (same as the server):
//
//	DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_SSLMODE
//
// The script is idempotent: it skips e-mails that already exist and
// skips duplicate applications (job+candidate pair).
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"recruitment-selection/internal/model"
)

// ── helpers ──────────────────────────────────────────────────────────────────

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func must(err error) {
	if err != nil {
		log.Fatalf("seed: %v", err)
	}
}

func hashPw(pw string) string {
	h, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	must(err)
	return string(h)
}

func ptr[T any](v T) *T { return &v }

// ── connect ───────────────────────────────────────────────────────────────────

func connect() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		env("DB_HOST", "localhost"),
		env("DB_PORT", "5432"),
		env("DB_USER", "postgres"),
		env("DB_PASSWORD", "postgres"),
		env("DB_NAME", "recruitment_selection"),
		env("DB_SSLMODE", "disable"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	must(err)
	return db
}

// ── upsert helpers ────────────────────────────────────────────────────────────

// upsertUser inserts a user or loads the existing one by e-mail.
func upsertUser(ctx context.Context, db *gorm.DB, email, name string, role model.UserRole) model.User {
	var u model.User
	err := db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if err == nil {
		return u
	}
	u = model.User{
		ID:           uuid.New(),
		Name:         name,
		Email:        email,
		PasswordHash: hashPw("senha123"),
		Role:         role,
	}
	must(db.WithContext(ctx).Create(&u).Error)
	fmt.Printf("  ✔ user %s (%s)\n", email, role)
	return u
}

// ── data definitions ─────────────────────────────────────────────────────────

type stageDef struct{ name string }

type jobDef struct {
	company     string
	title       string
	description string
	location    string
	salaryMin   float64
	salaryMax   float64
	status      model.JobStatus
	stages      []string
}

var jobs = []jobDef{
	{
		company:     "TechBR",
		title:       "Desenvolvedor Backend Go",
		description: "Desenvolvimento de microsserviços em Go, APIs REST e gRPC. Integração com Kafka e PostgreSQL. Experiência mínima de 3 anos.",
		location:    "São Paulo, SP (Híbrido)",
		salaryMin:   10000, salaryMax: 16000,
		status: model.JobStatusOpen,
		stages: []string{"Triagem", "Teste Técnico", "Entrevista Tech", "Entrevista RH", "Proposta"},
	},
	{
		company:     "TechBR",
		title:       "Desenvolvedor Frontend React",
		description: "Construção de SPAs com React e TypeScript. Experiência com design systems, testes e performance web.",
		location:    "Remoto",
		salaryMin:   8000, salaryMax: 13000,
		status: model.JobStatusOpen,
		stages: []string{"Triagem", "Desafio Frontend", "Entrevista Tech", "Entrevista Cultural"},
	},
	{
		company:     "Banco Digital",
		title:       "Analista de Dados Sênior",
		description: "Análise exploratória, modelagem preditiva e dashboards em Power BI e Metabase. Python e SQL avançados obrigatórios.",
		location:    "São Paulo, SP (Presencial)",
		salaryMin:   9000, salaryMax: 14000,
		status: model.JobStatusOpen,
		stages: []string{"Triagem", "Case Analítico", "Banca Técnica", "Entrevista Diretoria"},
	},
	{
		company:     "Banco Digital",
		title:       "Engenheiro de Segurança",
		description: "Gestão de vulnerabilidades, pentest e conformidade com LGPD e PCI-DSS. Certificação CISSP ou CEH desejável.",
		location:    "São Paulo, SP (Presencial)",
		salaryMin:   14000, salaryMax: 20000,
		status: model.JobStatusPaused,
		stages: []string{"Triagem", "Entrevista Técnica", "Banca de Segurança", "Proposta"},
	},
	{
		company:     "Startup XYZ",
		title:       "Full Stack Developer",
		description: "Startup de fintech busca full stack com Node.js + React. Ambiente ágil, equity disponível.",
		location:    "Remoto",
		salaryMin:   7000, salaryMax: 11000,
		status: model.JobStatusOpen,
		stages: []string{"Triagem", "Mini Projeto", "Pair Programming", "Fit Cultural"},
	},
	{
		company:     "Consultoria ABC",
		title:       "Gerente de Projetos TI",
		description: "Gestão de portfólio de projetos digitais, PMO, metodologias ágeis e certificação PMP.",
		location:    "Rio de Janeiro, RJ (Presencial)",
		salaryMin:   11000, salaryMax: 17000,
		status: model.JobStatusClosed,
		stages: []string{"Triagem", "Entrevista RH", "Entrevista Técnica", "Apresentação ao Cliente"},
	},
	{
		company:     "E-commerce Fast",
		title:       "DevOps / SRE Engineer",
		description: "Infraestrutura AWS, Kubernetes, observabilidade com Datadog. Suporte a plataforma de alta disponibilidade.",
		location:    "Remoto",
		salaryMin:   12000, salaryMax: 18000,
		status: model.JobStatusOpen,
		stages: []string{"Triagem", "Teste de Infraestrutura", "Entrevista SRE", "Proposta"},
	},
	{
		company:     "Saúde+",
		title:       "Product Manager",
		description: "PM para produto de telemedicina. Experiência com discovery, roadmap e KPIs de produto digital na área de saúde.",
		location:    "São Paulo, SP (Híbrido)",
		salaryMin:   13000, salaryMax: 19000,
		status: model.JobStatusCancelled,
		stages: []string{"Triagem", "Estudo de Caso", "Entrevista CEO", "Fit Time"},
	},
}

// ── main ──────────────────────────────────────────────────────────────────────

func main() {
	ctx := context.Background()
	db := connect()

	fmt.Println("\n=== SEED: users ===")

	recruiter := upsertUser(ctx, db, "recrutador@techbr.com", "Ana Recruiter", model.RoleRecruiter)

	candidateEmails := []struct{ email, name string }{
		{"joao.silva@email.com", "João Silva"},
		{"maria.santos@email.com", "Maria Santos"},
		{"pedro.alves@email.com", "Pedro Alves"},
		{"carla.ferreira@email.com", "Carla Ferreira"},
		{"lucas.oliveira@email.com", "Lucas Oliveira"},
		{"julia.costa@email.com", "Júlia Costa"},
		{"andre.lima@email.com", "André Lima"},
		{"fernanda.melo@email.com", "Fernanda Melo"},
		{"rafael.souza@email.com", "Rafael Souza"},
		{"beatriz.rocha@email.com", "Beatriz Rocha"},
	}

	candidates := make([]model.User, len(candidateEmails))
	for i, c := range candidateEmails {
		candidates[i] = upsertUser(ctx, db, c.email, c.name, model.RoleCandidate)
	}

	fmt.Println("\n=== SEED: jobs + stages ===")

	createdJobs := make([]model.Job, len(jobs))
	for i, jd := range jobs {
		var existing model.Job
		err := db.WithContext(ctx).
			Where("recruiter_id = ? AND title = ? AND company = ?", recruiter.ID, jd.title, jd.company).
			First(&existing).Error

		if err == nil {
			// reload with stages
			db.WithContext(ctx).
				Preload("Stages", func(db *gorm.DB) *gorm.DB { return db.Order("order_index ASC") }).
				First(&existing, "id = ?", existing.ID)
			createdJobs[i] = existing
			fmt.Printf("  ~ job already exists: %s / %s\n", jd.company, jd.title)
			continue
		}

		job := model.Job{
			ID:          uuid.New(),
			RecruiterID: recruiter.ID,
			Company:     jd.company,
			Title:       jd.title,
			Description: jd.description,
			Location:    jd.location,
			SalaryMin:   &jd.salaryMin,
			SalaryMax:   &jd.salaryMax,
			Status:      jd.status,
		}
		must(db.WithContext(ctx).Create(&job).Error)

		// Delete any stages the DB trigger may have auto-inserted (migration 005
		// created a trigger that inserts default English stages on every INSERT).
		// Migration 007 removes the trigger, but we guard here for safety.
		must(db.WithContext(ctx).
			Where("job_id = ?", job.ID).
			Delete(&model.JobStage{}).Error)

		stages := make([]model.JobStage, len(jd.stages))
		for si, sname := range jd.stages {
			stages[si] = model.JobStage{
				ID:         uuid.New(),
				JobID:      job.ID,
				Name:       sname,
				OrderIndex: si + 1,
			}
		}
		must(db.WithContext(ctx).Create(&stages).Error)
		job.Stages = stages
		createdJobs[i] = job
		fmt.Printf("  ✔ job %s / %s (%s)\n", jd.company, jd.title, jd.status)
	}

	fmt.Println("\n=== SEED: applications ===")

	type appSpec struct {
		jobIdx       int
		candidateIdx int
		status       model.ApplicationStatus
		stageIdx     int  // -1 = no stage yet (pending), >=0 = advance to this stage index
	}

	// Spread applications across all jobs with diverse stages
	specs := []appSpec{
		// Job 0 — Backend Go (open) — 4 candidates
		{0, 0, model.ApplicationStatusAccepted, 4},    // reached Proposta → accepted
		{0, 1, model.ApplicationStatusInProgress, 2},   // at Entrevista Tech
		{0, 2, model.ApplicationStatusInProgress, 1},   // at Teste Técnico
		{0, 3, model.ApplicationStatusRejected, 1},     // rejected at Teste Técnico

		// Job 1 — Frontend React (open) — 3 candidates
		{1, 4, model.ApplicationStatusInProgress, 2},   // at Entrevista Tech
		{1, 5, model.ApplicationStatusPending, -1},     // fresh, no stage
		{1, 6, model.ApplicationStatusWithdrawn, 0},    // withdrew at Triagem

		// Job 2 — Analista Dados (open) — 3 candidates
		{2, 0, model.ApplicationStatusInProgress, 2},   // at Banca Técnica
		{2, 7, model.ApplicationStatusInProgress, 1},   // at Case Analítico
		{2, 8, model.ApplicationStatusRejected, 0},     // rejected at Triagem

		// Job 3 — Segurança (paused) — 2 candidates
		{3, 1, model.ApplicationStatusInProgress, 2},   // at Banca de Segurança
		{3, 9, model.ApplicationStatusPending, -1},     // fresh

		// Job 4 — Full Stack XYZ (open) — 4 candidates
		{4, 2, model.ApplicationStatusAccepted, 3},     // reached Fit Cultural → accepted
		{4, 3, model.ApplicationStatusInProgress, 1},   // at Mini Projeto
		{4, 4, model.ApplicationStatusInProgress, 2},   // at Pair Programming
		{4, 5, model.ApplicationStatusRejected, 2},     // rejected at Pair Programming

		// Job 5 — Gerente Projetos (closed) — 2 candidates
		{5, 6, model.ApplicationStatusAccepted, 3},     // accepted
		{5, 7, model.ApplicationStatusRejected, 2},     // rejected

		// Job 6 — DevOps (open) — 3 candidates
		{6, 8, model.ApplicationStatusInProgress, 1},   // at Teste de Infraestrutura
		{6, 9, model.ApplicationStatusPending, -1},     // fresh
		{6, 0, model.ApplicationStatusInProgress, 2},   // at Entrevista SRE

		// Job 7 — Product Manager (cancelled) — 2 candidates
		{7, 1, model.ApplicationStatusRejected, 1},     // rejected
		{7, 2, model.ApplicationStatusWithdrawn, 0},    // withdrew
	}

	for _, sp := range specs {
		job := createdJobs[sp.jobIdx]
		candidate := candidates[sp.candidateIdx]

		// Skip if application already exists
		var existing model.Application
		err := db.WithContext(ctx).
			Where("job_id = ? AND candidate_id = ?", job.ID, candidate.ID).
			First(&existing).Error
		if err == nil {
			fmt.Printf("  ~ skip dup: %s → %s\n", candidate.Name, job.Title)
			continue
		}

		app := model.Application{
			ID:          uuid.New(),
			JobID:       job.ID,
			CandidateID: candidate.ID,
			Status:      sp.status,
			CoverLetter: fmt.Sprintf("Olá! Sou %s e tenho muito interesse nesta vaga de %s na %s. Possuo experiência relevante e estaria feliz em contribuir com a equipe.", candidate.Name, job.Title, job.Company),
			CVUrl:       fmt.Sprintf("/uploads/cv_%s.pdf", candidate.ID.String()),
			CreatedAt:   time.Now().Add(-time.Duration(sp.candidateIdx*12+sp.jobIdx*3) * time.Hour),
		}

		// Set stage if applicable
		if sp.stageIdx >= 0 && sp.stageIdx < len(job.Stages) {
			app.CurrentStageID = &job.Stages[sp.stageIdx].ID
		}

		must(db.WithContext(ctx).Create(&app).Error)
		stageLabel := "sem etapa"
		if sp.stageIdx >= 0 && sp.stageIdx < len(job.Stages) {
			stageLabel = job.Stages[sp.stageIdx].Name
		}
		fmt.Printf("  ✔ %s → %s [%s / %s]\n", candidate.Name, job.Title, sp.status, stageLabel)
	}

	fmt.Println("\n✅ Seed concluído!")
	fmt.Println("\nCredenciais de acesso:")
	fmt.Println("  Recrutador : recrutador@techbr.com / senha123")
	fmt.Println("  Candidatos : joao.silva@email.com ... beatriz.rocha@email.com / senha123")
}
