--
-- PostgreSQL database dump
--

-- Dumped from database version 12.9 (Ubuntu 12.9-0ubuntu0.20.04.1)
-- Dumped by pg_dump version 12.9 (Ubuntu 12.9-0ubuntu0.20.04.1)

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
-- Name: update_solution(); Type: FUNCTION; Schema: public; Owner: lk
--

CREATE FUNCTION public.update_solution() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
UPDATE solutions SET check_result = 5 WHERE task_id = OLD.id;
RETURN NEW;
END;
$$;


ALTER FUNCTION public.update_solution() OWNER TO lk;

--
-- Name: update_solution_status(bigint); Type: PROCEDURE; Schema: public; Owner: lk
--

CREATE PROCEDURE public.update_solution_status(i bigint)
    LANGUAGE sql
    AS $$
UPDATE solutions SET check_result = 5 WHERE id = i;
$$;


ALTER PROCEDURE public.update_solution_status(i bigint) OWNER TO lk;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: solutions; Type: TABLE; Schema: public; Owner: lk
--

CREATE TABLE public.solutions (
                                  id bigint NOT NULL,
                                  task_id bigint,
                                  check_result integer DEFAULT 0 NOT NULL,
                                  tests_passed integer NOT NULL,
                                  tests_total integer NOT NULL,
                                  received_date_time timestamp with time zone DEFAULT '2022-01-01 00:00:00+03'::timestamp with time zone NOT NULL,
                                  source_code text DEFAULT ''::text NOT NULL,
                                  uid bigint,
                                  check_time double precision DEFAULT 0.0 NOT NULL,
                                  check_message text DEFAULT ''::text NOT NULL,
                                  compile_time double precision DEFAULT 0.0 NOT NULL,
                                  checked_date_time timestamp with time zone DEFAULT '2022-01-01 00:00:00+03'::timestamp with time zone NOT NULL
);


ALTER TABLE public.solutions OWNER TO lk;

--
-- Name: solutions_id_seq; Type: SEQUENCE; Schema: public; Owner: lk
--

CREATE SEQUENCE public.solutions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.solutions_id_seq OWNER TO lk;

--
-- Name: solutions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: lk
--

ALTER SEQUENCE public.solutions_id_seq OWNED BY public.solutions.id;


--
-- Name: tasks; Type: TABLE; Schema: public; Owner: lk
--

CREATE TABLE public.tasks (
                              id bigint NOT NULL,
                              title text NOT NULL,
                              description text NOT NULL,
                              hints text,
                              input text NOT NULL,
                              output text NOT NULL,
                              test_amount integer DEFAULT 0 NOT NULL,
                              tests text NOT NULL,
                              creator bigint,
                              is_private boolean DEFAULT false NOT NULL,
                              code character varying(10) DEFAULT NULL::character varying,
                              date timestamp with time zone NOT NULL
);


ALTER TABLE public.tasks OWNER TO lk;

--
-- Name: tasks_done; Type: TABLE; Schema: public; Owner: lk
--

CREATE TABLE public.tasks_done (
                                   uid bigint,
                                   task_id bigint
);


ALTER TABLE public.tasks_done OWNER TO lk;

--
-- Name: tasks_id_seq; Type: SEQUENCE; Schema: public; Owner: lk
--

CREATE SEQUENCE public.tasks_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.tasks_id_seq OWNER TO lk;

--
-- Name: tasks_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: lk
--

ALTER SEQUENCE public.tasks_id_seq OWNED BY public.tasks.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: lk
--

CREATE TABLE public.users (
                              id bigint NOT NULL,
                              username character varying(60) NOT NULL,
                              password character varying(60) NOT NULL,
                              email character varying(100) NOT NULL,
                              fullname character varying(100) NOT NULL,
                              avatar_url text DEFAULT '/media/avatars/default.jpg'::text NOT NULL,
                              joined_date timestamp with time zone DEFAULT '2022-01-01 00:00:00+03'::timestamp with time zone NOT NULL,
                              is_admin boolean DEFAULT false NOT NULL,
                              verified boolean DEFAULT false NOT NULL
);


ALTER TABLE public.users OWNER TO lk;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: lk
--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO lk;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: lk
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: solutions id; Type: DEFAULT; Schema: public; Owner: lk
--

ALTER TABLE ONLY public.solutions ALTER COLUMN id SET DEFAULT nextval('public.solutions_id_seq'::regclass);


--
-- Name: tasks id; Type: DEFAULT; Schema: public; Owner: lk
--

ALTER TABLE ONLY public.tasks ALTER COLUMN id SET DEFAULT nextval('public.tasks_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: lk
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: solutions solutions_pkey; Type: CONSTRAINT; Schema: public; Owner: lk
--

ALTER TABLE ONLY public.solutions
    ADD CONSTRAINT solutions_pkey PRIMARY KEY (id);


--
-- Name: tasks tasks_pkey; Type: CONSTRAINT; Schema: public; Owner: lk
--

ALTER TABLE ONLY public.tasks
    ADD CONSTRAINT tasks_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: lk
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: lk
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: lk
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- Name: tasks update_solution; Type: TRIGGER; Schema: public; Owner: lk
--

CREATE TRIGGER update_solution AFTER UPDATE ON public.tasks FOR EACH ROW EXECUTE FUNCTION public.update_solution();


--
-- Name: solutions solutions_task_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: lk
--

ALTER TABLE ONLY public.solutions
    ADD CONSTRAINT solutions_task_id_fkey FOREIGN KEY (task_id) REFERENCES public.tasks(id) ON DELETE CASCADE;


--
-- Name: solutions solutions_uid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: lk
--

ALTER TABLE ONLY public.solutions
    ADD CONSTRAINT solutions_uid_fkey FOREIGN KEY (uid) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: tasks_done tasks_done_task_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: lk
--

ALTER TABLE ONLY public.tasks_done
    ADD CONSTRAINT tasks_done_task_id_fkey FOREIGN KEY (task_id) REFERENCES public.tasks(id) ON DELETE CASCADE;


--
-- Name: tasks_done tasks_done_uid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: lk
--

ALTER TABLE ONLY public.tasks_done
    ADD CONSTRAINT tasks_done_uid_fkey FOREIGN KEY (uid) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--