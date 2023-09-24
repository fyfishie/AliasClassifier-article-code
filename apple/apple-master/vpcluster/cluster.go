/*
 * @Author: fyfishie
 * @Date: 2023-03-27:08
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-05-08:21
 * @Description: :)
 * @email: fyfishie@outlook.com
 */
package vpcluster

import (
	"aliasParseMaster/idgen"
	"aliasParseMaster/lib"
	"errors"
	"strconv"
)

var ErrVPAddrAlreadyExist = errors.New("vp addr already exist")

type addVPCells struct {
	vp  lib.VP
	res chan (error)
}

type Cluster struct {
	//each slave machine have 10(or ant other nums)HP, when detected dead it lost one
	//if an slave machine's HP falls below to zero, it is thought to be lost and we delete it
	//map[vp.ID]*lib.VP
	aliveVPMap map[string]*lib.VP
	cellChan   chan (addVPCells)
	// allLocationReqChan chan (struct{})
	// allLocationChan    chan ([]lib.Location)
	allVPReqChan chan (struct{})

	/*
		attention!!! dont't return arrays of vp pointer
		incase vp list in vp cluster maybe changed during the time

		format:map[VP.ID]lib.VP
	*/
	allVPResChan chan (map[string]lib.VP)

	idGen *idgen.Generator
}

func NewCluster() *Cluster {
	return &Cluster{
		aliveVPMap: make(map[string]*lib.VP),
		cellChan:   make(chan addVPCells),
		// allLocationReqChan: make(chan struct{}),
		// allLocationChan:    make(chan []lib.Location),
		allVPReqChan: make(chan struct{}),
		allVPResChan: make(chan map[string]lib.VP),
		idGen:        idgen.NewIdGen().Run(),
	}
}

// inner used
func (c *Cluster) addVP(vp lib.VP) (err error) {
	_, ok := c.aliveVPMap[vp.UID]
	if ok {
		return ErrVPAddrAlreadyExist
	}
	c.aliveVPMap[vp.UID] = &vp
	return nil
}

// add a new vp, if there exists a same one, err!=nil, otherwise err=nil
// Cluster auto process thread safety problem and all you should do is a call of c.AddVP
func (c *Cluster) AddVP(vpUrl string, vpName string) (UID string, err error) {
	resChan := make(chan error)
	newID := c.idGen.NewID()
	UID = vpName + "_" + strconv.Itoa(newID)
	cell := addVPCells{
		vp: lib.VP{
			UID: UID,
			URL: vpUrl,
		},
		res: resChan,
	}
	c.cellChan <- cell
	return UID, <-resChan
}

// check whether a slave machine is alive or not
// func (c *Cluster) slaveAlive(vp lib.VP) bool {
// 	client := http.Client{
// 		Timeout: time.Minute,
// 	}
// 	req, err := http.NewRequest("GET", "http://"+addr.IP+":"+addr.Port+"/api/alive", nil)
// 	if err != nil {
// 		logrus.Error("failed to get request, err = " + err.Error())
// 		return false
// 	}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return false
// 	}
// 	return resp.StatusCode == http.StatusOK
// }

// check around which slave machine didn't not response and we minus one HP of its
// func (c *Cluster) checkAround() {
// 	for addr, vp := range c.aliveVPMap {
// 		if !c.slaveAlive(*vp) {
// 			c.aliveVPMap[addr].HP--
// 		} else {
// 			c.aliveVPMap[addr].HP = 10
// 		}
// 		if c.aliveVPMap[addr].HP <= 0 {
// 			logrus.Info("an alias vp lost, ip = " + addr.IP + " port = " + addr.Port)
// 			delete(c.aliveVPMap, addr)
// 		}
// 	}
// }

func (c *Cluster) Run() *Cluster {
	go func() {
		for {
			select {
			case cell := <-c.cellChan:
				cell.res <- c.addVP(cell.vp)
			case <-c.allVPReqChan:
				c.allVPResChan <- c.allVP()
			}
		}
	}()
	return c
}

// func (c *Cluster) allLocation() []lib.Location {
// 	res := []lib.Location{}
// 	for _, v := range c.aliveVPMap {
// 		res = append(res, v.Location)
// 	}
// 	return res
// }

// auto processes thread safety problem, just write cluster.AllLocation
// func (c *Cluster) AllLocation() []lib.Location {
// 	c.allLocationReqChan <- struct{}{}
// 	return <-c.allLocationChan
// }

func (c *Cluster) AllVP() map[string]lib.VP {
	c.allVPReqChan <- struct{}{}
	return <-c.allVPResChan
}
func (c *Cluster) allVP() map[string]lib.VP {
	res := map[string]lib.VP{}
	for _, v := range c.aliveVPMap {
		res[v.UID] = *v
	}
	return res
}
