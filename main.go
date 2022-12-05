package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/application-research/whypfs-core"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var (
	// ErrInvalidCid is returned when a cid is invalid
	DB *gorm.DB
	// OsSignal signal used to shutdown
	OsSignal chan os.Signal
)

//	Job that pulls a CID from estuary DB to local blockstore.
//	by default will have 2 workers which is really small compared to the size of content that we can pull from
//	each shuttles
func main() {

	shuttle := flag.String("shuttle", "shuttle-7.estuary.tech", "shuttle")
	numOfWorkers := flag.Int("numOfWorkers", 10, "numOfWorkers")
	fromDateRange := flag.String("fromDateRange", "2022-10-01", "fromDateRange")
	toDateRange := flag.String("toDateRange", "2023-12-31", "toDateRange")
	//limit := flag.Int("limit", 100, "limit") // just to limit the record it pulls
	//blockstore := flag.String("blockstore", ".whypfs", "blockstore") // if you don't indicate, then it'll create a new flatfs blockstore

	//	initialize whypfs-core
	whypfsPeer, err := whypfs.NewNode(whypfs.NewNodeParams{
		Ctx:       context.Background(),
		Datastore: whypfs.NewInMemoryDatastore(),
	})
	if err != nil {
		panic(err)
	}

	whypfsPeer.BootstrapPeers(whypfs.DefaultBootstrapPeers())
	// peer with Estuary shuttles and api node.

	// 	initialize database
	var errOnDb error
	DB, errOnDb = setupDB()
	if errOnDb != nil {
		panic(errOnDb)
	}

	// 	query database for all the CIDs.
	cids, errOnQuery := QueryAllCidsFromContentsWithoutDealTable(*shuttle, *fromDateRange, *toDateRange)
	if errOnQuery != nil {
		panic(errOnQuery)
	}

	// 	go routine worker (1000)
	var numJobs = len(cids)
	var perJob = numJobs / *numOfWorkers
	fmt.Println("numJobs", numJobs, "perJob", perJob)
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	//	 slice the cids into groups
	var startCtr = 0
	var endCtr = perJob

	//	build the workers and process each slice of cids
	for w := 1; w <= *numOfWorkers; w++ {
		var rangeOfCids = (cids)[startCtr:endCtr]
		go worker(w, whypfsPeer, rangeOfCids, jobs, results)
		startCtr = endCtr + 1
		endCtr = endCtr + perJob
	}

	//	return the jobs
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)

	//	return the results
	for a := 1; a <= numJobs; a++ {
		<-results
	}
}

type Cids []struct {
	Cid DbCID
}

// > Query all the cids from the contents table without the deal table
func QueryAllCidsFromContentsWithoutDealTable(shuttle string, fromDate string, toDate string) (Cids, error) {
	var cids Cids

	//select c.cid from contents as c where c.location = (select handle from shuttles as s where s.host like '%shuttle-6%') and created_at between '2022-10-01' and '2023-12-31';
	//select count(*) from contents as c where c.id not in (select content from content_deals as cd where cd.deal_id < 1) and c.location = (select handle from shuttles as s where s.host like '%shuttle-8%') and created_at between '2022-10-01' and '2023-12-31';
	err := DB.Raw("select c.cid as cid from contents as c where c.location = (select handle from shuttles as s where s.host = ?) and created_at between ? and ?", shuttle, fromDate, toDate).Scan(&cids).Error
	if err != nil {
		panic(err)
	}
	return cids, nil

}

// > This function sets up a database connection and returns a pointer to a gorm.DB object
func setupDB() (*gorm.DB, error) { // it's a pointer to a gorm.DB

	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()

	dbHost, okHost := viper.Get("DB_HOST").(string)
	dbUser, okUser := viper.Get("DB_USER").(string)
	dbPass, okPass := viper.Get("DB_PASS").(string)
	dbName, okName := viper.Get("DB_NAME").(string)
	dbPort, okPort := viper.Get("DB_PORT").(string)
	if !okHost || !okUser || !okPass || !okName || !okPort {
		panic("invalid database configuration")
	}

	dsn := "host=" + dbHost + " user=" + dbUser + " password=" + dbPass + " dbname=" + dbName + " port=" + dbPort + " sslmode=disable TimeZone=Asia/Shanghai"

	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return DB, nil
}

// 	Simple (but efficient) worker?
func worker(id int, peer *whypfs.Node, groupCid Cids, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		// IPLD get
		for _, cid := range groupCid {
			fmt.Println("worker", id, "started  job", j, "cid", cid)
			node, err := peer.Get(context.Background(), cid.Cid.CID) // pull!!
			if err != nil {
				fmt.Println("error... moving on")
			} else {
				fmt.Println("success", node.Cid())
			}
		}

		//	log the result somewhere? maybe a SQLite db.

		fmt.Println("worker", id, "finished job", j)
		results <- j * 2
	}
}
