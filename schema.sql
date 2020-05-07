--
-- PostgreSQL database dump
--

-- Dumped from database version 12.2 (Ubuntu 12.2-4)
-- Dumped by pg_dump version 12.2 (Ubuntu 12.2-4)

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
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: ballots; Type: TABLE; Schema: public; Owner: jwlarocque
--

CREATE TABLE public.ballots (
    ballot_id integer NOT NULL,
    question_id uuid NOT NULL,
    user_id character varying(80) NOT NULL
);


ALTER TABLE public.ballots OWNER TO jwlarocque;

--
-- Name: ballots_ballot_id_seq; Type: SEQUENCE; Schema: public; Owner: jwlarocque
--

ALTER TABLE public.ballots ALTER COLUMN ballot_id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.ballots_ballot_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: options; Type: TABLE; Schema: public; Owner: jwlarocque
--

CREATE TABLE public.options (
    option_id integer NOT NULL,
    question_id uuid NOT NULL,
    text character varying(128)
);


ALTER TABLE public.options OWNER TO jwlarocque;

--
-- Name: questions; Type: TABLE; Schema: public; Owner: jwlarocque
--

CREATE TABLE public.questions (
    name character varying(128),
    user_id character varying(80) NOT NULL,
    type integer NOT NULL,
    question_id uuid DEFAULT public.uuid_generate_v1() NOT NULL
);


ALTER TABLE public.questions OWNER TO jwlarocque;

--
-- Name: sessions; Type: TABLE; Schema: public; Owner: jwlarocque
--

CREATE TABLE public.sessions (
    session_id character varying(80) NOT NULL,
    user_id character varying(80) NOT NULL,
    created timestamp without time zone NOT NULL,
    expires timestamp without time zone NOT NULL
);


ALTER TABLE public.sessions OWNER TO jwlarocque;

--
-- Name: users; Type: TABLE; Schema: public; Owner: jwlarocque
--

CREATE TABLE public.users (
    user_id character varying(80) NOT NULL,
    email character varying(256)
);


ALTER TABLE public.users OWNER TO jwlarocque;

--
-- Name: votes; Type: TABLE; Schema: public; Owner: jwlarocque
--

CREATE TABLE public.votes (
    ballot_id integer NOT NULL,
    option_id integer NOT NULL,
    state integer NOT NULL
);


ALTER TABLE public.votes OWNER TO jwlarocque;

--
-- Name: ballots ballots_ballot_id_key; Type: CONSTRAINT; Schema: public; Owner: jwlarocque
--

ALTER TABLE ONLY public.ballots
    ADD CONSTRAINT ballots_ballot_id_key UNIQUE (ballot_id);


--
-- Name: ballots ballots_pkey; Type: CONSTRAINT; Schema: public; Owner: jwlarocque
--

ALTER TABLE ONLY public.ballots
    ADD CONSTRAINT ballots_pkey PRIMARY KEY (question_id, user_id);


--
-- Name: options options_pkey; Type: CONSTRAINT; Schema: public; Owner: jwlarocque
--

ALTER TABLE ONLY public.options
    ADD CONSTRAINT options_pkey PRIMARY KEY (option_id, question_id);


--
-- Name: questions questions_pk; Type: CONSTRAINT; Schema: public; Owner: jwlarocque
--

ALTER TABLE ONLY public.questions
    ADD CONSTRAINT questions_pk PRIMARY KEY (question_id);


--
-- Name: sessions session_pkey; Type: CONSTRAINT; Schema: public; Owner: jwlarocque
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT session_pkey PRIMARY KEY (user_id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: jwlarocque
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (user_id);


--
-- Name: votes votes_pkey; Type: CONSTRAINT; Schema: public; Owner: jwlarocque
--

ALTER TABLE ONLY public.votes
    ADD CONSTRAINT votes_pkey PRIMARY KEY (ballot_id, option_id);


--
-- Name: ballots ballots_question_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: jwlarocque
--

ALTER TABLE ONLY public.ballots
    ADD CONSTRAINT ballots_question_id_fkey FOREIGN KEY (question_id) REFERENCES public.questions(question_id);


--
-- Name: ballots ballots_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: jwlarocque
--

ALTER TABLE ONLY public.ballots
    ADD CONSTRAINT ballots_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(user_id);


--
-- Name: options options_question_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: jwlarocque
--

ALTER TABLE ONLY public.options
    ADD CONSTRAINT options_question_id_fkey FOREIGN KEY (question_id) REFERENCES public.questions(question_id);


--
-- Name: questions questions_user_id; Type: FK CONSTRAINT; Schema: public; Owner: jwlarocque
--

ALTER TABLE ONLY public.questions
    ADD CONSTRAINT questions_user_id FOREIGN KEY (user_id) REFERENCES public.users(user_id);


--
-- Name: sessions sessions_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: jwlarocque
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(user_id);


--
-- Name: votes votes_ballot_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: jwlarocque
--

ALTER TABLE ONLY public.votes
    ADD CONSTRAINT votes_ballot_id_fkey FOREIGN KEY (ballot_id) REFERENCES public.ballots(ballot_id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

