CREATE TABLE public.history (
	id serial NOT NULL,
	"timestamp" timestamp,
	CONSTRAINT history_pk PRIMARY KEY (id)

);

ALTER TABLE public.history OWNER TO postgres;


