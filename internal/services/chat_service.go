package services

import (
	"database/sql"
	"fmt"
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/internal/utils"
	"strings"
)

type ChatService struct {
	repo   repository.ChatRepository
	gemini *utils.GeminiClient
	db     *sql.DB
}

func NewChatService(
	r repository.ChatRepository,
	g *utils.GeminiClient,
	db *sql.DB,
) *ChatService {
	return &ChatService{
		repo:   r,
		gemini: g,
		db:     db,
	}
}

// GetDatabaseSchema: sql generate
func (s *ChatService) GetDatabaseSchema() (string, error) {
	query := `
		SELECT 
			TABLE_NAME,
			COLUMN_NAME,
			DATA_TYPE,
			COLUMN_KEY,
			IS_NULLABLE
		FROM INFORMATION_SCHEMA.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE()
		ORDER BY TABLE_NAME, ORDINAL_POSITION
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var schema strings.Builder
	schema.WriteString("DATABASE SCHEMA:\n\n")

	currentTable := ""
	for rows.Next() {
		var tableName, columnName, dataType, columnKey, isNullable string

		err := rows.Scan(&tableName, &columnName, &dataType, &columnKey, &isNullable)
		if err != nil {
			return "", err
		}

		if tableName != currentTable {
			if currentTable != "" {
				schema.WriteString("\n")
			}
			schema.WriteString(fmt.Sprintf("Table: %s\n", tableName))
			currentTable = tableName
		}

		nullInfo := ""
		if isNullable == "NO" {
			nullInfo = " NOT NULL"
		}

		keyInfo := ""
		if columnKey == "PRI" {
			keyInfo = " PRIMARY KEY"
		}

		schema.WriteString(fmt.Sprintf("  - %s: %s%s%s\n",
			columnName, dataType, nullInfo, keyInfo))
	}

	return schema.String(), nil
}

// GenerateSQLQuery - PURE LLM, NO INSTRUCTIONS
func (s *ChatService) GenerateSQLQuery(userQuestion string) (string, error) {
	schema, err := s.GetDatabaseSchema()
	if err != nil {
		return "", fmt.Errorf("failed to extract schema: %w", err)
	}

	// MINIMAL PROMPT - biarkan AI yang nalar sendiri
	prompt := fmt.Sprintf(`%s

User question (in Indonesian): %s

Generate the SQL query to answer this question. Return ONLY the SQL query without any explanation or markdown formatting.`, schema, userQuestion)

	sqlQuery, err := s.gemini.Generate(prompt)
	if err != nil {
		return "", err
	}

	// Clean up
	sqlQuery = strings.TrimSpace(sqlQuery)
	sqlQuery = strings.ReplaceAll(sqlQuery, "```sql", "")
	sqlQuery = strings.ReplaceAll(sqlQuery, "```", "")
	sqlQuery = strings.TrimSpace(sqlQuery)

	// Remove any trailing semicolon
	sqlQuery = strings.TrimSuffix(sqlQuery, ";")

	// Validasi SELECT only (security)
	upperQuery := strings.ToUpper(strings.TrimSpace(sqlQuery))
	if !strings.HasPrefix(upperQuery, "SELECT") {
		return "", fmt.Errorf("only SELECT queries allowed, got: %s", sqlQuery)
	}

	return sqlQuery, nil
}

// ExecuteQuery - Execute SQL
func (s *ChatService) ExecuteQuery(sqlQuery string) ([]map[string]interface{}, error) {
	rows, err := s.db.Query(sqlQuery)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}
		results = append(results, row)
	}

	return results, nil
}

// FormatQueryResults - AI format hasil jadi natural language
func (s *ChatService) FormatQueryResults(userQuestion string, results []map[string]interface{}) (string, error) {
	if len(results) == 0 {
		return "Tidak ada data yang ditemukan.", nil
	}

	var resultText strings.Builder

	for i, row := range results {
		resultText.WriteString(fmt.Sprintf("Row %d: ", i+1))
		parts := []string{}
		for key, val := range row {
			parts = append(parts, fmt.Sprintf("%s=%v", key, val))
		}
		resultText.WriteString(strings.Join(parts, ", "))
		resultText.WriteString("\n")
	}

	// MINIMAL PROMPT - biarkan AI yang tau cara format
	prompt := fmt.Sprintf(`User asked: "%s"

Database returned:
%s

Format this data as a natural, friendly response in Indonesian. Make it conversational and easy to read.`,
		userQuestion, resultText.String())

	answer, err := s.gemini.Generate(prompt)
	if err != nil {
		return "", err
	}

	return answer, nil
}

func (s *ChatService) Repo() repository.ChatRepository {
	return s.repo
}

func (s *ChatService) GetGemini() *utils.GeminiClient {
	return s.gemini
}
