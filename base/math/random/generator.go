package random

import (
)

type Generator struct {
	z0, z1, z2, z3 uint32
}

func Make() Generator {
	g := Generator{}
	g.Seed(0, 1, 2, 3)
	return g
}

func (g *Generator) Seed(s0, s1, s2, s3 uint32) {
	g.z0 = s0 | 128
	g.z1 = s1 | 128
	g.z2 = s2 | 128
	g.z3 = s3 | 128
}

func (g *Generator) advance_lfsr113() uint32 {
	var b uint32
	b  = ((g.z0 << 6) ^ g.z0) >> 13
	g.z0 = ((g.z0 & uint32(4294967294)) << 18) ^ b
	b  = ((g.z1 << 2) ^ g.z1) >> 27
	g.z1 = ((g.z1 & uint32(4294967288)) << 2) ^ b
	b  = ((g.z2 << 13) ^ g.z2) >> 21
	g.z2 = ((g.z2 & uint32(4294967280)) << 7) ^ b
	b  = ((g.z3 << 3) ^ g.z3) >> 12
	g.z3 = ((g.z3 & uint32(4294967168)) << 13) ^ b

	return g.z0 ^ g.z1 ^ g.z2 ^ g.z3
}

func (g *Generator) RandomFloat32() float32 {
	bits := g.advance_lfsr113()

	return 2.3283064365387e-10 * float32(bits)
}

/*
	//this is a simple lfsr113 RNG
	class RandomGenerator 
	{

	public:

		RandomGenerator(unsigned int s0 = 0, unsigned int s1 = 1, unsigned int s2 = 2, unsigned int s3 = 3) 
		{
			seed(s0, s1, s2, s3);
		}

		void seed(unsigned int s0, unsigned int s1, unsigned int s2, unsigned int s3) 
		{
			z0 = s0 | 128; // we need to make sure each value is > 127
			z1 = s1 | 128; // we need to make sure each value is > 127
			z2 = s2 | 128; // we need to make sure each value is > 127
			z3 = s3 | 128; // we need to make sure each value is > 127
		}


		unsigned int advance_lfsr113() 
		{
			unsigned int b;
			b  = ((z0 << 6) ^ z0) >> 13;
			z0 = ((z0 & 4294967294UL) << 18) ^ b;
			b  = ((z1 << 2) ^ z1) >> 27; 
			z1 = ((z1 & 4294967288UL) << 2) ^ b;
			b  = ((z2 << 13) ^ z2) >> 21;
			z2 = ((z2 & 4294967280UL) << 7) ^ b;
			b  = ((z3 << 3) ^ z3) >> 12;
			z3 = ((z3 & 4294967168UL) << 13) ^ b;

			return z0 ^ z1 ^ z2 ^ z3;
		}

		float randomNumber01() 
		{
			//return rand()/(float)RAND_MAX; //!!TODO
			unsigned int randomBits = advance_lfsr113();
			return 2.3283064365387e-10f * randomBits;
		}

	private:

		unsigned int z0;
		unsigned int z1;
		unsigned int z2;
		unsigned int z3;
	};

	*/