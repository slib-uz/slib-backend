package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"slib.uz/src/infrastructure/persistence/models"
)

type soatoRegion struct {
	SoatoID  int             `json:"soato_id"`
	NameUz   string          `json:"name_uz"`
	NameOz   string          `json:"name_oz"`
	NameRu   string          `json:"name_ru"`
	Children []soatoDistrict `json:"children"`
}

type soatoDistrict struct {
	SoatoID int    `json:"soato_id"`
	NameUz  string `json:"name_uz"`
	NameOz  string `json:"name_oz"`
	NameRu  string `json:"name_ru"`
}

func main() {
	soatoPath := flag.String("file", "assets/soato.json", "path to soato.json")
	flag.Parse()

	regions, err := readSoatoFile(*soatoPath)
	if err != nil {
		log.Fatalf("read soato file: %v", err)
	}

	db, err := openDB()
	if err != nil {
		log.Fatalf("open database: %v", err)
	}

	regionCount, districtCount, err := loadSoato(db, regions)
	if err != nil {
		log.Fatalf("load soato data: %v", err)
	}

	log.Printf("loaded %d regions and %d districts", regionCount, districtCount)
}

func readSoatoFile(path string) ([]soatoRegion, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	var regions []soatoRegion
	if err := json.Unmarshal(data, &regions); err != nil {
		return nil, err
	}

	return regions, nil
}

func openDB() (*gorm.DB, error) {
	_ = godotenv.Load("env/.env.local")
	_ = godotenv.Load("env/.env")

	required := []string{"DB_HOST", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_PORT", "DB_SSL_MODE"}
	for _, key := range required {
		if os.Getenv(key) == "" {
			return nil, fmt.Errorf("missing required env var: %s", key)
		}
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSL_MODE"),
	)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
}

func loadSoato(db *gorm.DB, regions []soatoRegion) (int, int, error) {
	regionCount := 0
	districtCount := 0

	err := db.Transaction(func(tx *gorm.DB) error {
		for _, region := range regions {
			regionModel := models.NewRegionModel(buildNameJSON(region.NameUz, region.NameOz, region.NameRu), region.SoatoID)
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "soato"}},
				DoUpdates: clause.AssignmentColumns([]string{"name", "updated_at"}),
			}).Create(regionModel).Error; err != nil {
				return fmt.Errorf("upsert region soato=%d: %w", region.SoatoID, err)
			}
			regionCount++
		}

		regionIDsBySoato, err := loadRegionIDsBySoato(tx)
		if err != nil {
			return err
		}

		for _, region := range regions {
			regionID, ok := regionIDsBySoato[region.SoatoID]
			if !ok {
				return fmt.Errorf("region not found after upsert: soato=%d", region.SoatoID)
			}

			for _, district := range region.Children {
				districtModel := models.NewDistrictModel(
					buildNameJSON(district.NameUz, district.NameOz, district.NameRu),
					district.SoatoID,
					regionID,
				)
				if err := tx.Clauses(clause.OnConflict{
					Columns:   []clause.Column{{Name: "soato"}},
					DoUpdates: clause.AssignmentColumns([]string{"name", "region_id", "updated_at"}),
				}).Create(districtModel).Error; err != nil {
					return fmt.Errorf("upsert district soato=%d: %w", district.SoatoID, err)
				}
				districtCount++
			}
		}

		return nil
	})
	if err != nil {
		return 0, 0, err
	}

	return regionCount, districtCount, nil
}

func loadRegionIDsBySoato(tx *gorm.DB) (map[int]uint, error) {
	var regionModels []models.RegionModel
	if err := tx.Select("id", "soato").Find(&regionModels).Error; err != nil {
		return nil, fmt.Errorf("load regions: %w", err)
	}

	regionIDsBySoato := make(map[int]uint, len(regionModels))
	for _, region := range regionModels {
		regionIDsBySoato[region.Soato] = region.ID
	}

	return regionIDsBySoato, nil
}

func buildNameJSON(uz, oz, ru string) datatypes.JSON {
	payload, _ := json.Marshal(map[string]string{
		"uz": uz,
		"oz": oz,
		"ru": ru,
	})
	return datatypes.JSON(payload)
}
