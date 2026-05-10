--
-- PostgreSQL database dump
--

\restrict GjkHaXpUl73nrrZFdwID9X8A4xNQzKL5s9IqTExHdaiZw634XgQ3b7hAd1cJ2I0

-- Dumped from database version 18.3
-- Dumped by pg_dump version 18.3

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: kasir; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.kasir (
    id_kasir bigint NOT NULL,
    no_telp text,
    tempat_lahir text,
    tanggal_lahir date,
    jenis_kelamin text,
    alamat text,
    pendidikan_terakhir text,
    nik character varying(16) NOT NULL,
    id_user bigint
);


ALTER TABLE public.kasir OWNER TO postgres;

--
-- Name: kasir_id_kasir_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.kasir_id_kasir_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.kasir_id_kasir_seq OWNER TO postgres;

--
-- Name: kasir_id_kasir_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.kasir_id_kasir_seq OWNED BY public.kasir.id_kasir;


--
-- Name: kurir; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.kurir (
    id_kurir bigint NOT NULL,
    no_telp character varying(15) NOT NULL,
    tempat_lahir text,
    tanggal_lahir date,
    jenis_kelamin text,
    alamat text,
    pendidikan_terakhir text,
    nik character varying(16) NOT NULL,
    id_user bigint
);


ALTER TABLE public.kurir OWNER TO postgres;

--
-- Name: kurir_id_kurir_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.kurir_id_kurir_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.kurir_id_kurir_seq OWNER TO postgres;

--
-- Name: kurir_id_kurir_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.kurir_id_kurir_seq OWNED BY public.kurir.id_kurir;


--
-- Name: kasir id_kasir; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.kasir ALTER COLUMN id_kasir SET DEFAULT nextval('public.kasir_id_kasir_seq'::regclass);


--
-- Name: kurir id_kurir; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.kurir ALTER COLUMN id_kurir SET DEFAULT nextval('public.kurir_id_kurir_seq'::regclass);


--
-- Data for Name: kasir; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.kasir (id_kasir, no_telp, tempat_lahir, tanggal_lahir, jenis_kelamin, alamat, pendidikan_terakhir, nik, id_user) FROM stdin;
\.


--
-- Data for Name: kurir; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.kurir (id_kurir, no_telp, tempat_lahir, tanggal_lahir, jenis_kelamin, alamat, pendidikan_terakhir, nik, id_user) FROM stdin;
\.


--
-- Name: kasir_id_kasir_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.kasir_id_kasir_seq', 1, false);


--
-- Name: kurir_id_kurir_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.kurir_id_kurir_seq', 1, false);


--
-- Name: kasir kasir_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.kasir
    ADD CONSTRAINT kasir_pkey PRIMARY KEY (id_kasir);


--
-- Name: kurir kurir_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.kurir
    ADD CONSTRAINT kurir_pkey PRIMARY KEY (id_kurir);


--
-- Name: kasir uni_kasir_id_user; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.kasir
    ADD CONSTRAINT uni_kasir_id_user UNIQUE (id_user);


--
-- Name: kasir uni_kasir_nik; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.kasir
    ADD CONSTRAINT uni_kasir_nik UNIQUE (nik);


--
-- Name: kurir uni_kurir_id_user; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.kurir
    ADD CONSTRAINT uni_kurir_id_user UNIQUE (id_user);


--
-- Name: kurir uni_kurir_nik; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.kurir
    ADD CONSTRAINT uni_kurir_nik UNIQUE (nik);


--
-- Name: kurir uni_kurir_no_telp; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.kurir
    ADD CONSTRAINT uni_kurir_no_telp UNIQUE (no_telp);


--
-- Name: kasir fk_kasir_user; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.kasir
    ADD CONSTRAINT fk_kasir_user FOREIGN KEY (id_user) REFERENCES public."user"(id_user);


--
-- Name: kurir fk_kurir_user; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.kurir
    ADD CONSTRAINT fk_kurir_user FOREIGN KEY (id_user) REFERENCES public."user"(id_user);


--
-- PostgreSQL database dump complete
--

\unrestrict GjkHaXpUl73nrrZFdwID9X8A4xNQzKL5s9IqTExHdaiZw634XgQ3b7hAd1cJ2I0

