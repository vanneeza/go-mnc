--
-- PostgreSQL database dump
--

-- Dumped from database version 15.2
-- Dumped by pg_dump version 15.2

-- Started on 2023-06-30 16:37:58

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
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
-- TOC entry 214 (class 1259 OID 28126)
-- Name: balance; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.balance (
    id character varying(255) NOT NULL,
    balance numeric(15,2) NOT NULL
);


ALTER TABLE public.balance OWNER TO postgres;

--
-- TOC entry 215 (class 1259 OID 28129)
-- Name: bank; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.bank (
    id character varying(255) NOT NULL,
    name character varying(100) NOT NULL,
    bank_account character varying(100) NOT NULL,
    branch character varying(100) NOT NULL,
    account_number bigint NOT NULL,
    merchant_id character varying(255) NOT NULL
);


ALTER TABLE public.bank OWNER TO postgres;

--
-- TOC entry 216 (class 1259 OID 28134)
-- Name: bank_admin; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.bank_admin (
    id character varying(255) NOT NULL,
    name character varying(100),
    bank_account character varying(100) NOT NULL,
    branch character varying(100) NOT NULL,
    account_number bigint NOT NULL
);


ALTER TABLE public.bank_admin OWNER TO postgres;

--
-- TOC entry 217 (class 1259 OID 28139)
-- Name: customer; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.customer (
    id character varying(255) NOT NULL,
    name character varying(100) NOT NULL,
    phone character varying(15) NOT NULL,
    address text NOT NULL,
    password character varying(100) NOT NULL,
    is_deleted boolean DEFAULT false,
    role character varying(10) DEFAULT 'customer'::character varying
);


ALTER TABLE public.customer OWNER TO postgres;

--
-- TOC entry 218 (class 1259 OID 28146)
-- Name: detail; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.detail (
    id character varying(255) NOT NULL,
    status character varying(50) NOT NULL,
    total_price numeric(15,2) NOT NULL,
    photo character varying(255) NOT NULL,
    bank_id character varying(255) NOT NULL
);


ALTER TABLE public.detail OWNER TO postgres;

--
-- TOC entry 219 (class 1259 OID 28151)
-- Name: log_user; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.log_user (
    id bigint NOT NULL,
    user_id character varying(255) NOT NULL,
    level character varying(10) NOT NULL,
    activity character varying(50) NOT NULL,
    description text NOT NULL,
    created_at timestamp with time zone NOT NULL
);


ALTER TABLE public.log_user OWNER TO postgres;

--
-- TOC entry 220 (class 1259 OID 28156)
-- Name: log_user_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.log_user_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.log_user_id_seq OWNER TO postgres;

--
-- TOC entry 3406 (class 0 OID 0)
-- Dependencies: 220
-- Name: log_user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.log_user_id_seq OWNED BY public.log_user.id;


--
-- TOC entry 221 (class 1259 OID 28157)
-- Name: merchant; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.merchant (
    id character varying(255) NOT NULL,
    name character varying(100) NOT NULL,
    phone character varying(15) NOT NULL,
    password character varying(255) NOT NULL,
    is_deleted boolean DEFAULT false,
    role character varying(10) DEFAULT 'merchant'::character varying
);


ALTER TABLE public.merchant OWNER TO postgres;

--
-- TOC entry 222 (class 1259 OID 28164)
-- Name: payment; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.payment (
    id character varying(255) NOT NULL,
    pay double precision NOT NULL,
    detail_id character varying(255) NOT NULL
);


ALTER TABLE public.payment OWNER TO postgres;

--
-- TOC entry 223 (class 1259 OID 28169)
-- Name: product; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.product (
    id character varying(255) NOT NULL,
    name character varying(100) NOT NULL,
    price numeric(15,2) NOT NULL,
    description text NOT NULL,
    is_deleted boolean DEFAULT false,
    merchant_id character varying(255) NOT NULL
);


ALTER TABLE public.product OWNER TO postgres;

--
-- TOC entry 224 (class 1259 OID 28175)
-- Name: tx_order; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.tx_order (
    id character varying(255) NOT NULL,
    qty smallint NOT NULL,
    product_id character varying(255) NOT NULL,
    customer_id character varying(255) NOT NULL,
    detail_id character varying(255) NOT NULL
);


ALTER TABLE public.tx_order OWNER TO postgres;

--
-- TOC entry 3211 (class 2604 OID 28180)
-- Name: log_user id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.log_user ALTER COLUMN id SET DEFAULT nextval('public.log_user_id_seq'::regclass);
--

SELECT pg_catalog.setval('public.log_user_id_seq', 1, false);


--
-- TOC entry 3216 (class 2606 OID 28182)
-- Name: balance balance_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.balance
    ADD CONSTRAINT balance_pkey PRIMARY KEY (id);


--
-- TOC entry 3222 (class 2606 OID 28184)
-- Name: bank_admin bank_admin_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.bank_admin
    ADD CONSTRAINT bank_admin_pkey PRIMARY KEY (id);


--
-- TOC entry 3218 (class 2606 OID 28186)
-- Name: bank bank_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.bank
    ADD CONSTRAINT bank_pkey PRIMARY KEY (id);


--
-- TOC entry 3224 (class 2606 OID 28188)
-- Name: customer customer_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.customer
    ADD CONSTRAINT customer_pkey PRIMARY KEY (id);


--
-- TOC entry 3228 (class 2606 OID 28190)
-- Name: detail detail_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.detail
    ADD CONSTRAINT detail_pkey PRIMARY KEY (id);


--
-- TOC entry 3230 (class 2606 OID 28192)
-- Name: log_user log_user_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.log_user
    ADD CONSTRAINT log_user_pkey PRIMARY KEY (id);


--
-- TOC entry 3232 (class 2606 OID 28194)
-- Name: merchant merchant_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.merchant
    ADD CONSTRAINT merchant_pkey PRIMARY KEY (id);


--
-- TOC entry 3240 (class 2606 OID 28196)
-- Name: tx_order order_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tx_order
    ADD CONSTRAINT order_pkey PRIMARY KEY (id);


--
-- TOC entry 3236 (class 2606 OID 28198)
-- Name: payment payment_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment
    ADD CONSTRAINT payment_pkey PRIMARY KEY (id);


--
-- TOC entry 3238 (class 2606 OID 28200)
-- Name: product product_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.product
    ADD CONSTRAINT product_pkey PRIMARY KEY (id);


--
-- TOC entry 3220 (class 2606 OID 28202)
-- Name: bank unique_account_number; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.bank
    ADD CONSTRAINT unique_account_number UNIQUE (account_number);


--
-- TOC entry 3234 (class 2606 OID 28204)
-- Name: merchant unique_phone; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.merchant
    ADD CONSTRAINT unique_phone UNIQUE (phone);


--
-- TOC entry 3226 (class 2606 OID 28206)
-- Name: customer unique_phone_customer; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.customer
    ADD CONSTRAINT unique_phone_customer UNIQUE (phone);


--
-- TOC entry 3245 (class 2606 OID 28207)
-- Name: tx_order f_product_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tx_order
    ADD CONSTRAINT f_product_id FOREIGN KEY (product_id) REFERENCES public.product(id) NOT VALID;


--
-- TOC entry 3242 (class 2606 OID 28212)
-- Name: detail fk_bank_admin_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.detail
    ADD CONSTRAINT fk_bank_admin_id FOREIGN KEY (bank_id) REFERENCES public.bank_admin(id) NOT VALID;


--
-- TOC entry 3246 (class 2606 OID 28217)
-- Name: tx_order fk_customer_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tx_order
    ADD CONSTRAINT fk_customer_id FOREIGN KEY (customer_id) REFERENCES public.customer(id);


--
-- TOC entry 3247 (class 2606 OID 28222)
-- Name: tx_order fk_detail_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tx_order
    ADD CONSTRAINT fk_detail_id FOREIGN KEY (detail_id) REFERENCES public.detail(id) NOT VALID;


--
-- TOC entry 3243 (class 2606 OID 28227)
-- Name: payment fk_detail_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment
    ADD CONSTRAINT fk_detail_id FOREIGN KEY (detail_id) REFERENCES public.detail(id);


--
-- TOC entry 3241 (class 2606 OID 28232)
-- Name: bank fk_merchant_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.bank
    ADD CONSTRAINT fk_merchant_id FOREIGN KEY (merchant_id) REFERENCES public.merchant(id);


--
-- TOC entry 3244 (class 2606 OID 28237)
-- Name: product fk_merchant_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.product
    ADD CONSTRAINT fk_merchant_id FOREIGN KEY (merchant_id) REFERENCES public.merchant(id) NOT VALID;


-- Completed on 2023-06-30 16:37:58

--
-- PostgreSQL database dump complete
--

