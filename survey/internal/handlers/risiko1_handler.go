package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"survey-backend/model"
	"survey-backend/repository"
	"survey-backend/service"
)

// Dependency wiring (singleton for in-memory demo)
// In production, inject these via constructor or DI framework.
var (
	ipTheftRepo   = repository.NewIPTheftRepository()
	progressRepo  = repository.NewProgressRepository()
	ipTheftSvc    = service.NewIPTheftService(ipTheftRepo, progressRepo)
	navigationSvc = service.NewNavigationService(progressRepo)
)

// POST /api/survey/risk/ip-theft/eligibility
// STEP 1 — Pertanyaan awal:
//   "Apakah perusahaan Anda berpotensi mengalami atau pernah mengalami
//    insiden pencurian Intellectual Property?"
// Request:
//   { "respondent_id": "R001", "has_experienced": true|false }
// Response:
//   { "next_step": "show_detail" }   ← jika Ya
//   { "next_step": "show_reason" }   ← jika Tidak

func SubmitIPTheftEligibility(w http.ResponseWriter, r *http.Request) {
	var req model.EligibilityRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := ipTheftSvc.ProcessEligibility(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, model.APIResponse{
		Success: true,
		Message: "Jawaban eligibilitas berhasil disimpan",
		Data:    result,
	})
}

// POST /api/survey/risk/ip-theft/reason
// STEP 2a — Alur "Tidak"
// Syarat: has_experienced = false
// Pertanyaan:
//   "Mengapa perusahaan Anda tidak berpotensi mengalami atau tidak pernah
//    mengalami insiden pencurian Intellectual Property?"
// Request:
//   { "respondent_id": "R001", "reason": "..." }
// Setelah ini → Risiko 1 selesai, tombol Berikutnya aktif

func SubmitIPTheftReason(w http.ResponseWriter, r *http.Request) {
	var req model.ReasonRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := ipTheftSvc.ProcessReason(req)
	if err != nil {
		writeError(w, resolveErrorStatus(err), err.Error())
		return
	}

	writeJSON(w, http.StatusOK, model.APIResponse{
		Success: true,
		Message: "Alasan berhasil disimpan. Risiko 1 selesai.",
		Data:    result,
	})
}

// POST /api/survey/risk/ip-theft/detail
// STEP 2b — Alur "Ya"
// Syarat: has_experienced = true
// Pertanyaan:
//   1. "Seberapa besar dampak dari pencurian Intellectual Property perusahaan?"
//      (matrix Reputasi / Operasional / Finansial / Hukum, nilai 1–4)
//   2. "Seberapa sering dalam setahun risiko pencurian IP berpotensi terjadi?"
//      (Kecil=1, Sedang=2, Besar=3, Sangat Besar=4)
// Request:
//   {
//     "respondent_id": "R001",
//     "impact": {
//       "reputation": 3, "operational": 2, "financial": 4, "legal": 2
//     },
//     "frequency": 3
//   }
// Response → next_step: "show_control"
//   UI harus menampilkan pertanyaan tindakan pengendalian (Step 2c)
//   sebelum tombol Berikutnya diaktifkan

func SubmitIPTheftDetail(w http.ResponseWriter, r *http.Request) {
	var req model.DetailRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := ipTheftSvc.ProcessDetail(req)
	if err != nil {
		writeError(w, resolveErrorStatus(err), err.Error())
		return
	}

	writeJSON(w, http.StatusOK, model.APIResponse{
		Success: true,
		Message: "Dampak dan frekuensi berhasil disimpan",
		Data:    result, // { next_step: "show_control" }
	})
}

// POST /api/survey/risk/ip-theft/control
// STEP 2c — Alur "Ya", sub-branching tindakan pengendalian
// Syarat: has_experienced = true AND step 2b sudah diisi
// Pertanyaan wajib:
//   "Apa perusahaan Anda telah memiliki tindakan pengendalian terhadap
//    risiko pencurian Intellectual Property?"
//   ● Ya  → wajib mengisi:
//            "Apa tindakan pengendalian yang telah dilakukan oleh perusahaan Anda
//             terhadap risiko pencurian Intellectual Property perusahaan?"
//   ● Tidak → tidak ada pertanyaan lanjutan, tombol Berikutnya langsung aktif
// Request (has_control = Ya):
//   {
//     "respondent_id": "R001",
//     "has_control": true,
//     "control_measures": "Enkripsi data, NDA karyawan, monitoring akses sistem"
//   }
// Request (has_control = Tidak):
//   {
//     "respondent_id": "R001",
//     "has_control": false
//   }
// Response → next_step: "finish"  (tombol Berikutnya aktif, Risiko 1 selesai)

func SubmitIPTheftControl(w http.ResponseWriter, r *http.Request) {
	var req model.ControlRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := ipTheftSvc.ProcessControl(req)
	if err != nil {
		writeError(w, resolveErrorStatus(err), err.Error())
		return
	}

	msg := "Tindakan pengendalian disimpan. Risiko 1 selesai."
	if !req.HasControl {
		msg = "Tidak ada tindakan pengendalian. Risiko 1 selesai."
	}

	writeJSON(w, http.StatusOK, model.APIResponse{
		Success: true,
		Message: msg,
		Data:    result, // { next_step: "finish" }
	})
}

// GET /api/survey/risk/ip-theft/{respondent_id}
func GetIPTheftResponse(w http.ResponseWriter, r *http.Request) {
	respondentID := r.PathValue("respondent_id")
	if respondentID == "" {
		writeError(w, http.StatusBadRequest, "respondent_id diperlukan")
		return
	}

	result, err := ipTheftSvc.GetResponse(respondentID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "Data tidak ditemukan")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, model.APIResponse{Success: true, Data: result})
}

// GET /api/survey/progress/{respondent_id}
func GetSurveyProgress(w http.ResponseWriter, r *http.Request) {
	respondentID := r.PathValue("respondent_id")
	if respondentID == "" {
		writeError(w, http.StatusBadRequest, "respondent_id diperlukan")
		return
	}

	progress, err := navigationSvc.GetProgress(respondentID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, model.APIResponse{Success: true, Data: progress})
}

// POST /api/survey/navigate
// Body: { "respondent_id": "R001", "direction": "next"|"previous", "current_risk": 1 }
func NavigateSurvey(w http.ResponseWriter, r *http.Request) {
	var req model.NavigateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	progress, err := navigationSvc.Navigate(req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, model.APIResponse{
		Success: true,
		Message: "Navigasi berhasil",
		Data:    progress,
	})
}

// Helpers
func decodeJSON(r *http.Request, dst interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(dst)
}

func writeJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, model.APIResponse{Success: false, Message: message})
}

func resolveErrorStatus(err error) int {
	if errors.Is(err, repository.ErrNotFound) {
		return http.StatusNotFound
	}
	return http.StatusBadRequest
}