--
-- PostgreSQL database dump
--

-- Dumped from database version 16.10 (Debian 16.10-1.pgdg13+1)
-- Dumped by pg_dump version 16.10 (Debian 16.10-1.pgdg13+1)

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

--
-- Name: event_type; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.event_type AS ENUM (
    'EARN',
    'SPEND'
);


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: accounts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.accounts (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    public_code text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    balance_points integer DEFAULT 0 NOT NULL,
    total_spend_money numeric(12,2) DEFAULT 0 NOT NULL,
    level_code text,
    CONSTRAINT chk_accounts_balance_nonnegative CHECK ((balance_points >= 0)),
    CONSTRAINT chk_accounts_total_spend_nonnegative CHECK ((total_spend_money >= (0)::numeric))
);


--
-- Name: accounts_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.accounts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: accounts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.accounts_id_seq OWNED BY public.accounts.id;


--
-- Name: events; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.events (
    id bigint NOT NULL,
    account_id bigint NOT NULL,
    type public.event_type NOT NULL,
    delta_points integer NOT NULL,
    balance_after integer NOT NULL,
    amount_money numeric(12,2),
    ts timestamp with time zone DEFAULT now() NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    ruleset_id bigint,
    actor_user_id bigint,
    CONSTRAINT chk_events_amount_money_nonnegative CHECK (((amount_money IS NULL) OR (amount_money >= (0)::numeric))),
    CONSTRAINT chk_events_balance_after_nonnegative CHECK ((balance_after >= 0))
);


--
-- Name: events_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.events_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: events_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.events_id_seq OWNED BY public.events.id;


--
-- Name: level_rules; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.level_rules (
    id bigint NOT NULL,
    ruleset_id bigint NOT NULL,
    level_code text NOT NULL,
    threshold_total_spend numeric(12,2) NOT NULL,
    percent_earn numeric(5,2) NOT NULL,
    CONSTRAINT chk_level_rules_percent_positive CHECK ((percent_earn > (0)::numeric)),
    CONSTRAINT chk_level_rules_threshold_nonnegative CHECK ((threshold_total_spend >= (0)::numeric))
);


--
-- Name: level_rules_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.level_rules_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: level_rules_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.level_rules_id_seq OWNED BY public.level_rules.id;


--
-- Name: operations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.operations (
    id bigint NOT NULL,
    operation_id text NOT NULL,
    account_id bigint NOT NULL,
    op_type public.event_type NOT NULL,
    request_json jsonb NOT NULL,
    response_json jsonb,
    http_status integer,
    event_id bigint,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: operations_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.operations_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: operations_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.operations_id_seq OWNED BY public.operations.id;


--
-- Name: roles; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.roles (
    code text NOT NULL
);


--
-- Name: ruleset; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ruleset (
    id bigint NOT NULL,
    effective_from timestamp with time zone NOT NULL,
    base_rub_per_point numeric(10,2) NOT NULL,
    created_by bigint,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT chk_ruleset_base_rub_per_point_positive CHECK ((base_rub_per_point > (0)::numeric))
);


--
-- Name: ruleset_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.ruleset_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: ruleset_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.ruleset_id_seq OWNED BY public.ruleset.id;


--
-- Name: user_roles; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_roles (
    user_id bigint NOT NULL,
    role_code text NOT NULL
);


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id bigint NOT NULL,
    phone text NOT NULL,
    password_hash text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    is_active boolean DEFAULT true NOT NULL
);


--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: accounts id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.accounts ALTER COLUMN id SET DEFAULT nextval('public.accounts_id_seq'::regclass);


--
-- Name: events id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.events ALTER COLUMN id SET DEFAULT nextval('public.events_id_seq'::regclass);


--
-- Name: level_rules id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.level_rules ALTER COLUMN id SET DEFAULT nextval('public.level_rules_id_seq'::regclass);


--
-- Name: operations id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.operations ALTER COLUMN id SET DEFAULT nextval('public.operations_id_seq'::regclass);


--
-- Name: ruleset id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ruleset ALTER COLUMN id SET DEFAULT nextval('public.ruleset_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: accounts accounts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.accounts
    ADD CONSTRAINT accounts_pkey PRIMARY KEY (id);


--
-- Name: events events_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.events
    ADD CONSTRAINT events_pkey PRIMARY KEY (id);


--
-- Name: level_rules level_rules_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.level_rules
    ADD CONSTRAINT level_rules_pkey PRIMARY KEY (id);


--
-- Name: operations operations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.operations
    ADD CONSTRAINT operations_pkey PRIMARY KEY (id);


--
-- Name: user_roles pk_user_roles; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT pk_user_roles PRIMARY KEY (user_id, role_code);


--
-- Name: roles roles_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_pkey PRIMARY KEY (code);


--
-- Name: ruleset ruleset_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ruleset
    ADD CONSTRAINT ruleset_pkey PRIMARY KEY (id);


--
-- Name: accounts uq_accounts_public_code; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.accounts
    ADD CONSTRAINT uq_accounts_public_code UNIQUE (public_code);


--
-- Name: accounts uq_accounts_user_id; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.accounts
    ADD CONSTRAINT uq_accounts_user_id UNIQUE (user_id);


--
-- Name: level_rules uq_level_rules_ruleset_level; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.level_rules
    ADD CONSTRAINT uq_level_rules_ruleset_level UNIQUE (ruleset_id, level_code);


--
-- Name: level_rules uq_level_rules_ruleset_threshold; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.level_rules
    ADD CONSTRAINT uq_level_rules_ruleset_threshold UNIQUE (ruleset_id, threshold_total_spend);


--
-- Name: operations uq_operations_idempotency; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.operations
    ADD CONSTRAINT uq_operations_idempotency UNIQUE (account_id, op_type, operation_id);


--
-- Name: ruleset uq_ruleset_effective_from; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ruleset
    ADD CONSTRAINT uq_ruleset_effective_from UNIQUE (effective_from);


--
-- Name: users uq_users_phone; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT uq_users_phone UNIQUE (phone);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_events_account_ts; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_events_account_ts ON public.events USING btree (account_id, ts);


--
-- Name: idx_events_actor_ts; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_events_actor_ts ON public.events USING btree (actor_user_id, ts);


--
-- Name: idx_operations_account_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_operations_account_created_at ON public.operations USING btree (account_id, created_at);


--
-- Name: accounts fk_accounts_user; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.accounts
    ADD CONSTRAINT fk_accounts_user FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE RESTRICT;


--
-- Name: events fk_events_account; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.events
    ADD CONSTRAINT fk_events_account FOREIGN KEY (account_id) REFERENCES public.accounts(id) ON DELETE RESTRICT;


--
-- Name: events fk_events_actor; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.events
    ADD CONSTRAINT fk_events_actor FOREIGN KEY (actor_user_id) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: events fk_events_ruleset; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.events
    ADD CONSTRAINT fk_events_ruleset FOREIGN KEY (ruleset_id) REFERENCES public.ruleset(id) ON DELETE SET NULL;


--
-- Name: level_rules fk_level_rules_ruleset; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.level_rules
    ADD CONSTRAINT fk_level_rules_ruleset FOREIGN KEY (ruleset_id) REFERENCES public.ruleset(id) ON DELETE CASCADE;


--
-- Name: operations fk_operations_account; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.operations
    ADD CONSTRAINT fk_operations_account FOREIGN KEY (account_id) REFERENCES public.accounts(id) ON DELETE RESTRICT;


--
-- Name: operations fk_operations_event; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.operations
    ADD CONSTRAINT fk_operations_event FOREIGN KEY (event_id) REFERENCES public.events(id) ON DELETE SET NULL;


--
-- Name: ruleset fk_ruleset_created_by; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ruleset
    ADD CONSTRAINT fk_ruleset_created_by FOREIGN KEY (created_by) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: user_roles fk_user_roles_role; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT fk_user_roles_role FOREIGN KEY (role_code) REFERENCES public.roles(code) ON DELETE RESTRICT;


--
-- Name: user_roles fk_user_roles_user; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT fk_user_roles_user FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

