/**
* Create background processes and pass out channels into them.
*
* @param {channel} out Where our random values will be written to
* @param {int} max All random numbers are < this number.
* @param {int} num_numbers How many random numbers to generate?
* @param {int} num_goroutines How many goroutines to use?
*/
func IntnBackground(out chan [][]uint64, max uint64, num_numbers int, 
	chunk_size uint64,
	num_goroutines int) {
  
  // Poisson implements the Poisson distribution, a discrete probability distribution
// that expresses the probability of a given number of events occurring in a fixed
// interval.
// The poisson distribution has density function:
//  f(k) = λ^k / k! e^(-λ)
// For more information, see https://en.wikipedia.org/wiki/Poisson_distribution.
type Poisson struct {
	// Lambda is the average number of events in an interval.
	// Lambda must be greater than 0.
	Lambda float64

	Src rand.Source
}

// CDF computes the value of the cumulative distribution function at x.
func (p Poisson) CDF(x float64) float64 {
	if x < 0 {
		return 0
	}
	return mathext.GammaIncRegComp(math.Floor(x+1), p.Lambda)
}

// ExKurtosis returns the excess kurtosis of the distribution.
func (p Poisson) ExKurtosis() float64 {
	return 1 / p.Lambda
}

// LogProb computes the natural logarithm of the value of the probability
// density function at x.
func (p Poisson) LogProb(x float64) float64 {
	if x < 0 || math.Floor(x) != x {
		return math.Inf(-1)
	}
	lg, _ := math.Lgamma(math.Floor(x) + 1)
	return x*math.Log(p.Lambda) - p.Lambda - lg
}

// Mean returns the mean of the probability distribution.
func (p Poisson) Mean() float64 {
	return p.Lambda
}

// NumParameters returns the number of parameters in the distribution.
func (Poisson) NumParameters() int {
	return 1
}

// Prob computes the value of the probability density function at x.
func (p Poisson) Prob(x float64) float64 {
	return math.Exp(p.LogProb(x))
}

// Rand returns a random sample drawn from the distribution.
func (p Poisson) Rand() float64 {
	// NUMERICAL RECIPES IN C: THE ART OF SCIENTIFIC COMPUTING (ISBN 0-521-43108-5)
	// p. 294
	// <http://www.aip.de/groups/soe/local/numres/bookcpdf/c7-3.pdf>

	rnd := rand.ExpFloat64
	var rng *rand.Rand
	if p.Src != nil {
		rng = rand.New(p.Src)
		rnd = rng.ExpFloat64
	}

	if p.Lambda < 10.0 {
		// Use direct method.
		var em float64
		t := 0.0
		for {
			t += rnd()
			if t >= p.Lambda {
				break
			}
			em++
		}
		return em
	}
	// Use rejection method.
	rnd = rand.Float64
	if rng != nil {
		rnd = rng.Float64
	}
	sq := math.Sqrt(2.0 * p.Lambda)
	alxm := math.Log(p.Lambda)
	lg, _ := math.Lgamma(p.Lambda + 1)
	g := p.Lambda*alxm - lg
	for {
		var em, y float64
		for {
			y = math.Tan(math.Pi * rnd())
			em = sq*y + p.Lambda
			if em >= 0 {
				break
			}
		}
		em = math.Floor(em)
		lg, _ = math.Lgamma(em + 1)
		t := 0.9 * (1.0 + y*y) * math.Exp(em*alxm-lg-g)
		if rnd() <= t {
			return em
		}
	}
}

// Skewness returns the skewness of the distribution.
func (p Poisson) Skewness() float64 {
	return 1 / math.Sqrt(p.Lambda)
}

// StdDev returns the standard deviation of the probability distribution.
func (p Poisson) StdDev() float64 {
	return math.Sqrt(p.Variance())
}

// Survival returns the survival function (complementary CDF) at x.
func (p Poisson) Survival(x float64) float64 {
	return 1 - p.CDF(x)
}

// Variance returns the variance of the probability distribution.
func (p Poisson) Variance() float64 {
	return p.Lambda
}

	log.Info("Starting background number generation...")

	in := make(chan []uint64)

	//
	// Create a number of background processes.
	//
	for i := 0; i < num_goroutines; i++ {
		random_struct := random_struct{false}
		go random_struct.intNChannel(in, out)
	}

	//
	// Now stuff our input channel with all of our requests.
	// We'll do it in chunks so as not to destroy our CPUs.
	// See my blog post at http://www.dmuth.org/node/1414/multi-core-cpu-performance-google-go
	// for a further explanation.
	//
	num_left := uint64(num_numbers)

	for {
		if (num_left < chunk_size) {
			chunk_size = num_left
		}

		num_left -= chunk_size

		var args []uint64
		args = append(args, max, chunk_size)
		in <- args

		if (num_left <= 0) {
			break
		}

	}

} // End of IntnBackground()


