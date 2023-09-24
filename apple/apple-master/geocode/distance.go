/*
 * @Author: fyfishie
 * @Date: 2023-03-27:10
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-05-08:19
 * @Description: :)
 * @email: fyfishie@outlook.com
 */
package geocode

// calculates distance of two points
// type Ruler struct {
// 	cache               *freecache.Cache
// 	cacheLength         int
// 	cacheInnerMapLength int
// }

// func (r *Ruler) DisTance(locationA, locationB lib.Location) float64 {
// 	d := r.queryCache(locationA, locationB)
// 	if d != -1 {
// 		return d
// 	}
// 	d = calculate(locationA, locationB)
// 	r.cache.Set(gpss2Key(locationA.GPS, locationB.GPS), utils.Float642Bytes(d), -1)
// 	return d
// }

// // TODO:we need cache, certainly!
// // -1 represents that no cache hitted
// func (r *Ruler) queryCache(locationA, locationB lib.Location) float64 {
// 	key := gpss2Key(locationA.GPS, locationB.GPS)
// 	got, err := r.cache.Get(key)
// 	if err != nil {
// 		return -1
// 	}
// 	return utils.Bytes2Float64(got)
// }

// func calculate(locationA, locationB lib.Location) float64 {
// 	return math.Sqrt(math.Pow(locationA.GPS.Latitude-locationB.GPS.Latitude, 2) + math.Pow(locationA.GPS.Longitude-locationB.GPS.Longitude, 2))
// }
// func NewRuler() *Ruler {
// 	return &Ruler{
// 		cache:               freecache.NewCache(100 * 1024 * 1024),
// 		cacheLength:         128,
// 		cacheInnerMapLength: 128,
// 	}

// }

// // just contact four result of float642Bytes and this construct is fun?...
// func gpss2Key(gps1, gps2 lib.GPSData) []byte {
// 	return append(
// 		append(
// 			append(
// 				utils.Float642Bytes(gps1.Latitude),
// 				utils.Float642Bytes(gps2.Longitude)...),
// 			utils.Float642Bytes(gps2.Latitude)...),
// 		utils.Float642Bytes(gps2.Longitude)...,
// 	)
// }
// func (r *Ruler) Nearest(target lib.Location, list []lib.Location) lib.Location {
// 	minDistance := math.MaxFloat64
// 	minDisIndex := 0
// 	for i, item := range list {
// 		d := r.DisTance(target, item)
// 		if d < minDistance {
// 			minDistance = d
// 			minDisIndex = i
// 		}
// 	}
// 	return list[minDisIndex]
// }
