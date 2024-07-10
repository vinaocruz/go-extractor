CREATE TABLE public.negotiations (
	id SERIAL NOT NULL,
	closedat VARCHAR NOT NULL,
	transationat DATE NOT NULL,
	ticketcode VARCHAR NOT NULL,
	price FLOAT4 NOT NULL,
	quantity INTEGER NOT NULL,
	CONSTRAINT negotiations_pk PRIMARY KEY (id)
);