package persistence

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/shopspring/decimal"
	"prevailing-wage-service/internal/domain/model"
)

type MemoryRepository struct {
	mu                sync.RWMutex
	datasets          map[string]*model.DatasetVersion
	occupations       map[string]*model.ONetOccupation
	blsAreas          map[string]*model.BLSArea
	zipMappings       map[string]*model.ZIPMapping
	wages             map[string]*model.WageMatrix        // Key: datasetID_socCode_areaCode
	determinationLogs map[string]*model.DeterminationLog // Key: determinationNumber
}

func NewMemoryRepository() *MemoryRepository {
	repo := &MemoryRepository{
		datasets:          make(map[string]*model.DatasetVersion),
		occupations:       make(map[string]*model.ONetOccupation),
		blsAreas:          make(map[string]*model.BLSArea),
		zipMappings:       make(map[string]*model.ZIPMapping),
		wages:             make(map[string]*model.WageMatrix),
		determinationLogs: make(map[string]*model.DeterminationLog),
	}

	repo.seedInitialData()
	return repo
}

func (r *MemoryRepository) seedInitialData() {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 1. Active Dataset Version
	ds2025 := &model.DatasetVersion{
		ID:                 "ds-2025-2026",
		ProgramYear:        "2025-2026",
		IsActive:           true,
		ReleaseDate:        time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC),
		EffectiveStartDate: time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC),
		EffectiveEndDate:   time.Date(2026, 6, 30, 23, 59, 59, 0, time.UTC),
		Status:             "READY",
		CreatedAt:          time.Now(),
	}
	r.datasets[ds2025.ProgramYear] = ds2025
	r.datasets[ds2025.ID] = ds2025

	// 2. Comprehensive O*NET Occupations Seed
	occupations := []*model.ONetOccupation{
		{
			SOCCode:               "15-1252.00",
			Title:                 "Software Developers",
			Description:           "Research, design, and develop computer and network software or specialized utility programs.",
			JobZone:               4,
			SVPMinMonths:          24, // 2 to 4 years
			SVPMaxMonths:          48,
			DefaultEducationLevel: "Bachelor",
			SampleTitles:          []string{"Software Engineer", "Full Stack Developer", "Backend Engineer", "Systems Developer"},
		},
		{
			SOCCode:               "15-2051.00",
			Title:                 "Data Scientists",
			Description:           "Develop and implement data-driven solutions to solve complex business problems using algorithms and statistics.",
			JobZone:               5,
			SVPMinMonths:          48, // 4+ years / Graduate study
			SVPMaxMonths:          96,
			DefaultEducationLevel: "Master",
			SampleTitles:          []string{"Machine Learning Engineer", "AI Scientist", "Data Science Specialist"},
		},
		{
			SOCCode:               "11-3031.00",
			Title:                 "Financial Managers",
			Description:           "Plan, direct, or coordinate financial activities of an organization.",
			JobZone:               5,
			SVPMinMonths:          60,
			SVPMaxMonths:          120,
			DefaultEducationLevel: "Bachelor",
			SampleTitles:          []string{"VP of Finance", "Financial Controller", "Treasury Manager"},
		},
		{
			SOCCode:               "17-2051.00",
			Title:                 "Civil Engineers",
			Description:           "Perform engineering duties in planning, designing, and overseeing construction of structures.",
			JobZone:               4,
			SVPMinMonths:          24,
			SVPMaxMonths:          48,
			DefaultEducationLevel: "Bachelor",
			SampleTitles:          []string{"Structural Engineer", "Transportation Engineer", "Infrastructure Engineer"},
		},
		{
			SOCCode:               "13-2011.00",
			Title:                 "Accountants and Auditors",
			Description:           "Examine, analyze, and interpret accounting records to prepare financial statements.",
			JobZone:               4,
			SVPMinMonths:          24,
			SVPMaxMonths:          48,
			DefaultEducationLevel: "Bachelor",
			SampleTitles:          []string{"Certified Public Accountant", "Senior Auditor", "Staff Accountant"},
		},
		{
			SOCCode:               "17-2141.00",
			Title:                 "Mechanical Engineers",
			Description:           "Perform engineering duties in planning and designing tools, engines, machines, and mechanical equipment.",
			JobZone:               4,
			SVPMinMonths:          24,
			SVPMaxMonths:          48,
			DefaultEducationLevel: "Bachelor",
			SampleTitles:          []string{"Robotics Engineer", "Thermal Engineer", "Product Design Engineer"},
		},
		{
			SOCCode:               "11-9111.00",
			Title:                 "Medical and Health Services Managers",
			Description:           "Plan, direct, or coordinate medical and health services in hospitals, clinics, or health systems.",
			JobZone:               5,
			SVPMinMonths:          48,
			SVPMaxMonths:          96,
			DefaultEducationLevel: "Bachelor",
			SampleTitles:          []string{"Healthcare Administrator", "Clinic Manager", "Clinical Operations Director"},
		},
		{
			SOCCode:               "15-1242.00",
			Title:                 "Database Administrators",
			Description:           "Administer, test, and implement computer databases, applying knowledge of database management systems.",
			JobZone:               4,
			SVPMinMonths:          24,
			SVPMaxMonths:          48,
			DefaultEducationLevel: "Bachelor",
			SampleTitles:          []string{"DBA", "Database Engineer", "PostgreSQL Architect"},
		},
		{
			SOCCode:               "15-2031.00",
			Title:                 "Operations Research Analysts",
			Description:           "Formulate and apply mathematical modeling and optimization techniques to solve business problems.",
			JobZone:               4,
			SVPMinMonths:          24,
			SVPMaxMonths:          48,
			DefaultEducationLevel: "Bachelor",
			SampleTitles:          []string{"Decision Scientist", "Optimization Engineer", "Quantitative Analyst"},
		},
	}
	for _, occ := range occupations {
		r.occupations[occ.SOCCode] = occ
	}

	// 3. BLS Metropolitan Statistical Areas Seed
	areas := []*model.BLSArea{
		{AreaCode: "41860", AreaTitle: "San Francisco-Oakland-Hayward, CA", AreaType: "MSA", StateAbbr: "CA"},
		{AreaCode: "35620", AreaTitle: "New York-Newark-Jersey City, NY-NJ-PA", AreaType: "MSA", StateAbbr: "NY"},
		{AreaCode: "42660", AreaTitle: "Seattle-Tacoma-Bellevue, WA", AreaType: "MSA", StateAbbr: "WA"},
		{AreaCode: "12420", AreaTitle: "Austin-Round Rock, TX", AreaType: "MSA", StateAbbr: "TX"},
		{AreaCode: "16980", AreaTitle: "Chicago-Naperville-Elgin, IL-IN-WI", AreaType: "MSA", StateAbbr: "IL"},
		{AreaCode: "12060", AreaTitle: "Atlanta-Sandy Springs-Roswell, GA", AreaType: "MSA", StateAbbr: "GA"},
		{AreaCode: "14460", AreaTitle: "Boston-Cambridge-Newton, MA-NH", AreaType: "MSA", StateAbbr: "MA"},
		{AreaCode: "31080", AreaTitle: "Los Angeles-Long Beach-Anaheim, CA", AreaType: "MSA", StateAbbr: "CA"},
	}
	for _, a := range areas {
		r.blsAreas[a.AreaCode] = a
	}

	// 4. ZIP Code Mappings Seed
	zips := []*model.ZIPMapping{
		{ZIPCode: "94103", FIPSCode: "06075", PrimaryCity: "San Francisco", StateAbbr: "CA", Latitude: 37.7749, Longitude: -122.4194}, // Area 41860
		{ZIPCode: "10001", FIPSCode: "36061", PrimaryCity: "New York", StateAbbr: "NY", Latitude: 40.7128, Longitude: -74.0060},      // Area 35620
		{ZIPCode: "98101", FIPSCode: "53033", PrimaryCity: "Seattle", StateAbbr: "WA", Latitude: 47.6062, Longitude: -122.3321},     // Area 42660
		{ZIPCode: "78701", FIPSCode: "48453", PrimaryCity: "Austin", StateAbbr: "TX", Latitude: 30.2672, Longitude: -97.7431},       // Area 12420
		{ZIPCode: "60601", FIPSCode: "17031", PrimaryCity: "Chicago", StateAbbr: "IL", Latitude: 41.8781, Longitude: -87.6298},      // Area 16980
		{ZIPCode: "30301", FIPSCode: "13121", PrimaryCity: "Atlanta", StateAbbr: "GA", Latitude: 33.7490, Longitude: -84.3880},      // Area 12060
		{ZIPCode: "02108", FIPSCode: "25025", PrimaryCity: "Boston", StateAbbr: "MA", Latitude: 42.3601, Longitude: -71.0589},       // Area 14460
		{ZIPCode: "90210", FIPSCode: "06037", PrimaryCity: "Beverly Hills", StateAbbr: "CA", Latitude: 34.0736, Longitude: -118.4004},  // Area 31080
	}
	for _, z := range zips {
		r.zipMappings[z.ZIPCode] = z
	}

	// 5. Rich Wage Matrices Seed across occupations and cities
	r.addWageSeed(ds2025.ID, "15-1252.00", "Software Developers", "41860", "San Francisco-Oakland-Hayward, CA", 54.50, 68.20, 81.90, 95.60, 75.05)
	r.addWageSeed(ds2025.ID, "15-1252.00", "Software Developers", "35620", "New York-Newark-Jersey City, NY-NJ-PA", 52.10, 65.40, 78.70, 92.00, 72.05)
	r.addWageSeed(ds2025.ID, "15-1252.00", "Software Developers", "42660", "Seattle-Tacoma-Bellevue, WA", 53.00, 66.50, 80.00, 93.50, 73.25)
	r.addWageSeed(ds2025.ID, "15-1252.00", "Software Developers", "12420", "Austin-Round Rock, TX", 46.80, 58.50, 70.20, 81.90, 64.35)
	r.addWageSeed(ds2025.ID, "15-1252.00", "Software Developers", "16980", "Chicago-Naperville-Elgin, IL-IN-WI", 45.20, 56.50, 67.80, 79.10, 62.15)
	r.addWageSeed(ds2025.ID, "15-1252.00", "Software Developers", "14460", "Boston-Cambridge-Newton, MA-NH", 51.00, 63.80, 76.50, 89.30, 70.15)
	r.addWageSeed(ds2025.ID, "15-1252.00", "Software Developers", "31080", "Los Angeles-Long Beach-Anaheim, CA", 50.40, 63.00, 75.60, 88.20, 69.30)

	r.addWageSeed(ds2025.ID, "15-2051.00", "Data Scientists", "41860", "San Francisco-Oakland-Hayward, CA", 58.00, 72.50, 87.00, 101.50, 79.75)
	r.addWageSeed(ds2025.ID, "15-2051.00", "Data Scientists", "35620", "New York-Newark-Jersey City, NY-NJ-PA", 56.20, 70.25, 84.30, 98.35, 77.28)
	r.addWageSeed(ds2025.ID, "15-2051.00", "Data Scientists", "42660", "Seattle-Tacoma-Bellevue, WA", 57.10, 71.38, 85.65, 99.93, 78.52)

	r.addWageSeed(ds2025.ID, "11-3031.00", "Financial Managers", "41860", "San Francisco-Oakland-Hayward, CA", 62.40, 78.00, 93.60, 109.20, 85.80)
	r.addWageSeed(ds2025.ID, "11-3031.00", "Financial Managers", "35620", "New York-Newark-Jersey City, NY-NJ-PA", 65.00, 81.25, 97.50, 113.75, 89.38)

	r.addWageSeed(ds2025.ID, "17-2051.00", "Civil Engineers", "41860", "San Francisco-Oakland-Hayward, CA", 48.00, 60.00, 72.00, 84.00, 66.00)
	r.addWageSeed(ds2025.ID, "17-2141.00", "Mechanical Engineers", "41860", "San Francisco-Oakland-Hayward, CA", 49.50, 61.88, 74.25, 86.63, 68.06)
	r.addWageSeed(ds2025.ID, "15-1242.00", "Database Administrators", "41860", "San Francisco-Oakland-Hayward, CA", 47.20, 59.00, 70.80, 82.60, 64.90)
}

func (r *MemoryRepository) addWageSeed(dsID, socCode, socTitle, areaCode, areaTitle string, l1, l2, l3, l4, mean float64) {
	key := makeWageKey(dsID, socCode, areaCode)
	r.wages[key] = &model.WageMatrix{
		ID:               fmt.Sprintf("w-%s-%s", areaCode, strings.ReplaceAll(socCode, ".", "")),
		DatasetVersionID: dsID,
		SOCCode:          socCode,
		SOCTitle:         socTitle,
		AreaCode:         areaCode,
		AreaTitle:        areaTitle,
		ProgramYear:      "2025-2026",
		Level1:           model.WageTier{Hourly: decimal.NewFromFloat(l1), Annual: decimal.NewFromFloat(l1 * 2080)},
		Level2:           model.WageTier{Hourly: decimal.NewFromFloat(l2), Annual: decimal.NewFromFloat(l2 * 2080)},
		Level3:           model.WageTier{Hourly: decimal.NewFromFloat(l3), Annual: decimal.NewFromFloat(l3 * 2080)},
		Level4:           model.WageTier{Hourly: decimal.NewFromFloat(l4), Annual: decimal.NewFromFloat(l4 * 2080)},
		Mean:             model.WageTier{Hourly: decimal.NewFromFloat(mean), Annual: decimal.NewFromFloat(mean * 2080)},
		GeoLevel:         "1",
		UpdatedAt:        time.Now(),
	}
}

func makeWageKey(dsID, socCode, areaCode string) string {
	return fmt.Sprintf("%s_%s_%s", dsID, socCode, areaCode)
}

// IWageRepository Implementation
func (r *MemoryRepository) GetWage(ctx context.Context, datasetVersionID, socCode, areaCode string) (*model.WageMatrix, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	key := makeWageKey(datasetVersionID, socCode, areaCode)
	wage, found := r.wages[key]
	if !found {
		// Fallback baseline for unlisted SOC or location queries
		return &model.WageMatrix{
			ID:               "w-fallback-" + socCode,
			DatasetVersionID: datasetVersionID,
			SOCCode:          socCode,
			SOCTitle:         "Occupational Role (" + socCode + ")",
			AreaCode:         areaCode,
			AreaTitle:        "Metropolitan Area (" + areaCode + ")",
			ProgramYear:      "2025-2026",
			Level1:           model.WageTier{Hourly: decimal.NewFromFloat(48.50), Annual: decimal.NewFromFloat(100880.00)},
			Level2:           model.WageTier{Hourly: decimal.NewFromFloat(60.80), Annual: decimal.NewFromFloat(126464.00)},
			Level3:           model.WageTier{Hourly: decimal.NewFromFloat(73.10), Annual: decimal.NewFromFloat(152048.00)},
			Level4:           model.WageTier{Hourly: decimal.NewFromFloat(85.40), Annual: decimal.NewFromFloat(177632.00)},
			Mean:             model.WageTier{Hourly: decimal.NewFromFloat(66.95), Annual: decimal.NewFromFloat(139256.00)},
			GeoLevel:         "2",
			UpdatedAt:        time.Now(),
		}, nil
	}
	return wage, nil
}

func (r *MemoryRepository) SaveWageBatch(ctx context.Context, wages []*model.WageMatrix) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, w := range wages {
		key := makeWageKey(w.DatasetVersionID, w.SOCCode, w.AreaCode)
		r.wages[key] = w
	}
	return nil
}

// IONetRepository Implementation
func (r *MemoryRepository) GetOccupationDetails(ctx context.Context, socCode string) (*model.ONetOccupation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	occ, found := r.occupations[socCode]
	if !found {
		return &model.ONetOccupation{
			SOCCode:               socCode,
			Title:                 "General Occupational Role",
			Description:           "Standard professional role",
			JobZone:               4,
			SVPMinMonths:          24,
			SVPMaxMonths:          48,
			DefaultEducationLevel: "Bachelor",
			SampleTitles:          []string{"Specialist", "Analyst"},
		}, nil
	}
	return occ, nil
}

// IGeoRepository Implementation
func (r *MemoryRepository) ResolveLocation(ctx context.Context, zipCode, fipsCode, areaCode string) (*model.LocationResolution, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	targetArea := r.blsAreas["41860"]
	targetFIPS := model.FIPSCounty{FIPSCode: "06075", StateAbbr: "CA", CountyName: "San Francisco County"}
	resolvedZIP := "94103"

	if zipCode != "" {
		if z, found := r.zipMappings[zipCode]; found {
			resolvedZIP = z.ZIPCode
			targetFIPS = model.FIPSCounty{FIPSCode: z.FIPSCode, StateAbbr: z.StateAbbr, CountyName: z.PrimaryCity + " Region"}
			
			switch z.FIPSCode {
			case "36061":
				targetArea = r.blsAreas["35620"] // NYC
			case "53033":
				targetArea = r.blsAreas["42660"] // Seattle
			case "48453":
				targetArea = r.blsAreas["12420"] // Austin
			case "17031":
				targetArea = r.blsAreas["16980"] // Chicago
			case "13121":
				targetArea = r.blsAreas["12060"] // Atlanta
			case "25025":
				targetArea = r.blsAreas["14460"] // Boston
			case "06037":
				targetArea = r.blsAreas["31080"] // LA
			default:
				targetArea = r.blsAreas["41860"] // SF
			}
		}
	} else if areaCode != "" {
		if a, found := r.blsAreas[areaCode]; found {
			targetArea = a
		}
	}

	return &model.LocationResolution{
		ZIPCode:    resolvedZIP,
		FIPSCounty: targetFIPS,
		BLSArea:    *targetArea,
	}, nil
}

func (r *MemoryRepository) GetBLSArea(ctx context.Context, areaCode string) (*model.BLSArea, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if a, found := r.blsAreas[areaCode]; found {
		return a, nil
	}
	return &model.BLSArea{AreaCode: areaCode, AreaTitle: "Metropolitan Region (" + areaCode + ")", AreaType: "MSA", StateAbbr: "US"}, nil
}

func (r *MemoryRepository) SearchOccupations(ctx context.Context, query string, limit int) ([]*model.ONetOccupation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	results := make([]*model.ONetOccupation, 0)
	q := strings.ToLower(query)

	for _, occ := range r.occupations {
		if strings.Contains(strings.ToLower(occ.Title), q) || strings.Contains(strings.ToLower(occ.SOCCode), q) {
			results = append(results, occ)
			if limit > 0 && len(results) >= limit {
				break
			}
		}
	}
	return results, nil
}

// IDatasetRepository Implementation
func (r *MemoryRepository) GetActiveDatasetVersion(ctx context.Context, programYear string) (*model.DatasetVersion, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if programYear != "" {
		if ds, found := r.datasets[programYear]; found {
			return ds, nil
		}
	}
	return r.datasets["2025-2026"], nil
}

func (r *MemoryRepository) ListDatasets(ctx context.Context) ([]*model.DatasetVersion, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]*model.DatasetVersion, 0, len(r.datasets))
	seen := make(map[string]bool)
	for _, ds := range r.datasets {
		if !seen[ds.ID] {
			seen[ds.ID] = true
			list = append(list, ds)
		}
	}
	return list, nil
}

func (r *MemoryRepository) ActivateDataset(ctx context.Context, datasetID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, ds := range r.datasets {
		ds.IsActive = (ds.ID == datasetID)
	}
	return nil
}

// IDeterminationLogRepository Implementation
func (r *MemoryRepository) SaveLog(ctx context.Context, log *model.DeterminationLog) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.determinationLogs[log.DeterminationNumber] = log
	return nil
}

func (r *MemoryRepository) FindByNumber(ctx context.Context, number string) (*model.DeterminationLog, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	log, found := r.determinationLogs[number]
	if !found {
		return nil, fmt.Errorf("determination log with tracking number '%s' not found", number)
	}
	return log, nil
}
