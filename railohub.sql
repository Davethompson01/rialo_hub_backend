--
-- PostgreSQL database dump
--

\restrict VLpe72oHAgemW72KrEdbAwA9fjsSBbR6136fY7ovS1PwcKB0MmOab4eluMl8Aqi

-- Dumped from database version 18.4 (Debian 18.4-1.pgdg13+1)
-- Dumped by pg_dump version 18.4 (Debian 18.4-1.pgdg13+1)

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

--
-- Name: status; Type: TYPE; Schema: public; Owner: dave
--

CREATE TYPE public.status AS ENUM (
    'Ongoing',
    'Accepted',
    'Completed'
);


ALTER TYPE public.status OWNER TO dave;

--
-- Name: task_application_status; Type: TYPE; Schema: public; Owner: dave
--

CREATE TYPE public.task_application_status AS ENUM (
    'Ongoing',
    'Accepted',
    'Rejected'
);


ALTER TYPE public.task_application_status OWNER TO dave;

--
-- Name: user_role; Type: TYPE; Schema: public; Owner: dave
--

CREATE TYPE public.user_role AS ENUM (
    'Writer',
    'Developer',
    'Artist',
    'Moderators',
    '3dArtist'
);


ALTER TYPE public.user_role OWNER TO dave;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: task_application; Type: TABLE; Schema: public; Owner: dave
--

CREATE TABLE public.task_application (
    task_application_id integer NOT NULL,
    task_id integer,
    employee_id integer,
    employer_id integer,
    skills text,
    status public.task_application_status,
    applied_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.task_application OWNER TO dave;

--
-- Name: task_application_task_application_id_seq; Type: SEQUENCE; Schema: public; Owner: dave
--

ALTER TABLE public.task_application ALTER COLUMN task_application_id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.task_application_task_application_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: tasks; Type: TABLE; Schema: public; Owner: dave
--

CREATE TABLE public.tasks (
    task_id integer NOT NULL,
    user_id integer,
    title text NOT NULL,
    description text NOT NULL,
    reward integer NOT NULL,
    status public.status,
    deadline timestamp without time zone CONSTRAINT tasks_dealine_not_null NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.tasks OWNER TO dave;

--
-- Name: tasks_task_id_seq; Type: SEQUENCE; Schema: public; Owner: dave
--

ALTER TABLE public.tasks ALTER COLUMN task_id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.tasks_task_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: users; Type: TABLE; Schema: public; Owner: dave
--

CREATE TABLE public.users (
    user_id integer NOT NULL,
    discord_username text NOT NULL,
    username text NOT NULL,
    role public.user_role,
    password text NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    profile_pics text,
    reputation text DEFAULT 0
);


ALTER TABLE public.users OWNER TO dave;

--
-- Name: users_user_id_seq; Type: SEQUENCE; Schema: public; Owner: dave
--

ALTER TABLE public.users ALTER COLUMN user_id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.users_user_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Data for Name: task_application; Type: TABLE DATA; Schema: public; Owner: dave
--

COPY public.task_application (task_application_id, task_id, employee_id, employer_id, skills, status, applied_at) FROM stdin;
14	12	1	1	Golang, PostgreSQL, RabbitMQ, REST APIs	Ongoing	2026-07-29 22:02:41.165997
15	12	2	1	Golang, PostgreSQL, RabbitMQ, REST APIs	Ongoing	2026-07-29 22:06:39.335632
17	11	2	1	Golang, PostgreSQL, RabbitMQ, REST APIs	Ongoing	2026-07-29 22:09:36.206071
\.


--
-- Data for Name: tasks; Type: TABLE DATA; Schema: public; Owner: dave
--

COPY public.tasks (task_id, user_id, title, description, reward, status, deadline, created_at) FROM stdin;
9	1	Develop Authentication Module	Implement secure JWT login and registration endpoints for the API.	500	Ongoing	2026-12-31 23:59:59	2026-07-29 19:23:18.080541
10	1	Develn Module	 registration endpoints for the API.	100	Ongoing	2026-12-31 23:59:59	2026-07-29 19:23:49.567788
11	2	 Module Develn	  endpoints for the Develn  API.	100	Ongoing	2026-12-31 23:59:59	2026-07-29 19:35:20.319933
12	2	Develop Authentication Module	 Authentication Module  endpoints for the Develn  API.	100	Ongoing	2026-12-31 23:59:59	2026-07-29 19:35:44.159831
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: dave
--

COPY public.users (user_id, discord_username, username, role, password, created_at, profile_pics, reputation) FROM stdin;
1	Dave	dave	Writer	$2a$10$mwFIwVIpdUJDzsBhs5tgM.UOjomzH0FI46kDhEKP5nGByhED8531m	2026-07-28 02:18:15.157239	\N	0
2	HayKeenz	Ayorinda	Moderators	$2a$10$anZ59noLzUrNuYV7WnQfGOov2Bu7w/pks6kyTKEoLP3/592WpGnba	2026-07-29 01:37:15.030501	\N	0
\.


--
-- Name: task_application_task_application_id_seq; Type: SEQUENCE SET; Schema: public; Owner: dave
--

SELECT pg_catalog.setval('public.task_application_task_application_id_seq', 17, true);


--
-- Name: tasks_task_id_seq; Type: SEQUENCE SET; Schema: public; Owner: dave
--

SELECT pg_catalog.setval('public.tasks_task_id_seq', 12, true);


--
-- Name: users_user_id_seq; Type: SEQUENCE SET; Schema: public; Owner: dave
--

SELECT pg_catalog.setval('public.users_user_id_seq', 2, true);


--
-- Name: task_application task_application_pkey; Type: CONSTRAINT; Schema: public; Owner: dave
--

ALTER TABLE ONLY public.task_application
    ADD CONSTRAINT task_application_pkey PRIMARY KEY (task_application_id);


--
-- Name: tasks tasks_pkey; Type: CONSTRAINT; Schema: public; Owner: dave
--

ALTER TABLE ONLY public.tasks
    ADD CONSTRAINT tasks_pkey PRIMARY KEY (task_id);


--
-- Name: users users_discord_username_key; Type: CONSTRAINT; Schema: public; Owner: dave
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_discord_username_key UNIQUE (discord_username);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: dave
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (user_id);


--
-- Name: task_application task_application_employee_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dave
--

ALTER TABLE ONLY public.task_application
    ADD CONSTRAINT task_application_employee_id_fkey FOREIGN KEY (employee_id) REFERENCES public.users(user_id);


--
-- Name: task_application task_application_employer_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dave
--

ALTER TABLE ONLY public.task_application
    ADD CONSTRAINT task_application_employer_id_fkey FOREIGN KEY (employer_id) REFERENCES public.users(user_id);


--
-- Name: task_application task_application_task_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dave
--

ALTER TABLE ONLY public.task_application
    ADD CONSTRAINT task_application_task_id_fkey FOREIGN KEY (task_id) REFERENCES public.tasks(task_id);


--
-- Name: tasks tasks_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dave
--

ALTER TABLE ONLY public.tasks
    ADD CONSTRAINT tasks_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(user_id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

\unrestrict VLpe72oHAgemW72KrEdbAwA9fjsSBbR6136fY7ovS1PwcKB0MmOab4eluMl8Aqi

