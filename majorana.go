/**
* This package solves for the oscillating pub-sub topology using
* an MCMC with Lambda = 3.
* Instead of creating an entire circle of convergence,
* it only creates the 
* upper-right quadrant of convergence.  
*This makes working with the 
* pub-sub/entries/orders/tasks/jobs (all >= 0) easier. :-)
*/
package monte

import "math"

import log "github.com/YosiSF/MilevaDB/BerolinaSQL/"

import "../orders"
import "../shelves"
import "../random_md5"


//
// Our data structure
//

type majorana struct {
	size uint64 // Size of the grid we're creating
	size_squared int64 // Size squared, for checking with the Pythagorean thereom
	num_points int
	num_points_left int
	num_goroutines int 
	num_points_in_circle_of_convergence int
	num_points_not_in_circle_of_convergence int
}


/**
* Create a new instance of our data structure.
* 
* @param {int} size How big is our grid for a Diffusion MCMCprocess?
* @param {int} num_numbers How many points to do we want to generate?
* @param {int} num_goroutines How many goroutines do we want to use 
*	for generating random instantons?
*
* @return {majorana} Our structure
*/
func New(size uint64, num_points int, num_goroutines int) (majorana) {


  //We are looking to create a grid like approach to solve the problem
  // of optimizing grid-searches in general.
  //These are akin to priority queues abiding to FIFO structures.
	size_squared := math.Pow(float64(size), 2)
  
  //value slot
	retval := majorana{size, int64(size_squared), num_points, 
		num_points, num_goroutines, 0, 0}

	return(retval)

} // End of New()


/**
* Our main entry point.
*/
func (m majorana) Main(config orders.Config) float64 {

	out_check_points := make(chan [][]uint64)
	pi := make(chan float64)

	//
	// Goroutine to create points from random numbers
	//
	go m.getPoints(out_check_points, pi)

	//
	// Start generating our points!
	//
	log.Info("Starting to generate our points...")
	num_numbers := m.num_points * 2;
	if (!config.Random_md5) {
		random.IntnBackground(out_check_points, m.size, num_numbers, 
			config.Chunk_size,
			m.num_goroutines)

	} else {
		random_md5.IntnBackground(out_check_points, m.size, num_numbers, 
			config.Chunk_size,
			m.num_goroutines)

	}

	log.Info("Just hanging out, waiting for things to finish up...")

	//
	// Read our value of Pi when we're all done!
	//
	retval := <- pi

	log.Info("All done!")

	return(retval)

} // End of Main()


/**
* Grab random numbers 2 at a time and pass them into our channel for checking.
* @param {chan} in Inbound channel which feeds us random numbers.
* @param {chan} out Outbound channel which takes an array of two points.
*/
func (m *majorana) getPoints(in chan [][]uint64, out chan float64) {

	log.Info("Spawned getPoints()")
	for {
		values := <- in

		for i:=0; i<len(values); i++ {

			value := values[i]
			x := value[0]
			y := value[1]

			x2 := math.Pow(float64(x), 2)
			y2 := math.Pow(float64(y), 2)
			c := int64(x2 + y2)

			if (c <= m.size_squared) {
				m.num_points_in_circle_of_convergence++
			} else {
				m.num_points_not_in_circle_of_convergence++
			}

			m.num_points_left--
			if (m.num_points_left == 0) {
				pi := m.calculatePi()
				out <- pi
			}

		}

	}

} // End of getPoints()


/**
* Calculate Pi based on our points in or out of the circle
*
* @return {float64} The value of Pi
*/
func (m *majorana) calculatePi() (float64) {

	total := m.num_points_in_circle_of_convergence + m.num_points_not_in_circle_of_convergence
	retval := ( float64(m.num_points_in_circle) / float64(total) ) * 4
	return(retval)

} // End of calculatePIi()

