package masterdata

import (
	"github.com/jszwec/csvutil"
	"io"
	"log"
	"os"
)

type MasterData struct {
	Item                *ItemManager
	Monster             *MonsterManager
	MonsterEnhanceTable *MonsterEnhanceTableManager
}

var master *MasterData

func init() {
	dummyItem, err := dummyItemData()
	if err != nil {
		log.Fatal(err)
	}
	dummyMonster, err := dummyMonsterData()
	if err != nil {
		log.Fatal(err)
	}
	dummyMonsterEnhanceTable, err := dummyMonsterEnhanceTableData()
	if err != nil {
		log.Fatal(err)
	}

	// 本来であればS3などから読み込む
	master = &MasterData{
		Item:                dummyItem,
		Monster:             dummyMonster,
		MonsterEnhanceTable: dummyMonsterEnhanceTable,
	}
	log.Println("success set masterdata")
}

func Master() *MasterData {
	return master
}

func dummyItemData() (*ItemManager, error) {
	b, err := readCsv("./pkg/masterdata/csv/item.csv")
	if err != nil {
		return nil, err
	}

	var items Items
	if err = csvutil.Unmarshal(b, &items); err != nil {
		return nil, err
	}

	return &ItemManager{
		items:   items,
		itemMap: items.ToMap(),
	}, nil
}

func dummyMonsterData() (*MonsterManager, error) {
	b, err := readCsv("./pkg/masterdata/csv/monster.csv")
	if err != nil {
		return nil, err
	}

	var monsters Monsters
	if err = csvutil.Unmarshal(b, &monsters); err != nil {
		return nil, err
	}

	return &MonsterManager{
		monsters:   monsters,
		monsterMap: monsters.ToMap(),
	}, nil
}

func dummyMonsterEnhanceTableData() (*MonsterEnhanceTableManager, error) {
	b, err := readCsv("./pkg/masterdata/csv/monster_enhance_table.csv")
	if err != nil {
		return nil, err
	}

	var monsterEnhanceTables MonsterEnhanceTables
	if err = csvutil.Unmarshal(b, &monsterEnhanceTables); err != nil {
		return nil, err
	}

	return &MonsterEnhanceTableManager{
		monsterEnhanceTables:   monsterEnhanceTables,
		monsterEnhanceTableMap: monsterEnhanceTables.ToMap(),
	}, nil
}

func readCsv(filepath string) ([]byte, error) {
	csvFile, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()

	b, err := io.ReadAll(csvFile)
	if err != nil {
		return nil, err
	}
	return b, nil
}
