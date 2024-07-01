--
-- PostgreSQL database dump
--

-- Dumped from database version 16.2
-- Dumped by pg_dump version 16.2

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
-- Name: api_token_type; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.api_token_type AS ENUM (
    'SYNC_PHOTO'
);


--
-- Name: auth_role; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.auth_role AS ENUM (
    'ADMIN',
    'USER'
);


--
-- Name: auth_status; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.auth_status AS ENUM (
    'NOT_ACTIVATED',
    'SENT_INVITE',
    'ACTIVATED',
    'BLOCKED'
);


--
-- Name: code_type; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.code_type AS ENUM (
    'ACTIVATE_INVITE',
    'ACTIVATE_REGISTRATION'
);


--
-- Name: photo_extension; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.photo_extension AS ENUM (
    'JPEG',
    'PNG'
);


--
-- Name: photo_status; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.photo_status AS ENUM (
    'ACTIVE',
    'NOT_VALID'
);


--
-- Name: processing_type; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.processing_type AS ENUM (
    'EXIF_DATA',
    'META_DATA',
    'CATALOG_TAGS',
    'META_TAGS',
    'PHOTO_VECTOR',
    'SIMILAR_COEFFICIENT',
    'PHOTO_GROUP',
    'PHOTO_PREVIEW'
);


--
-- Name: refresh_token_status; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.refresh_token_status AS ENUM (
    'ACTIVE',
    'REVOKED',
    'EXPIRED',
    'LOGOUT'
);


SET default_table_access_method = heap;

--
-- Name: api_tokens; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.api_tokens (
    id uuid NOT NULL,
    person_id uuid NOT NULL,
    caption text NOT NULL,
    token text NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    expired_at timestamp with time zone,
    type public.api_token_type NOT NULL,
    CONSTRAINT api_tokens_caption_check CHECK ((length(caption) <= 128)),
    CONSTRAINT api_tokens_caption_check1 CHECK ((length(caption) <= 32))
);


--
-- Name: auth; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.auth (
    person_id uuid NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    email text NOT NULL,
    password_hash bytea NOT NULL,
    status public.auth_status NOT NULL,
    role public.auth_role NOT NULL,
    CONSTRAINT auth_email_check CHECK ((length(email) <= 1024))
);


--
-- Name: codes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.codes (
    code text NOT NULL,
    person_id uuid NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    active boolean NOT NULL,
    type public.code_type NOT NULL
);


--
-- Name: coefficients_similar_photos; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.coefficients_similar_photos (
    photo_id1 uuid NOT NULL,
    photo_id2 uuid NOT NULL,
    coefficient double precision NOT NULL
);


--
-- Name: exif_photo_data; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.exif_photo_data (
    photo_id uuid NOT NULL,
    data jsonb
);


--
-- Name: goose_db_version; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.goose_db_version (
    id integer NOT NULL,
    version_id bigint NOT NULL,
    is_applied boolean NOT NULL,
    tstamp timestamp without time zone DEFAULT now()
);


--
-- Name: goose_db_version_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.goose_db_version_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: goose_db_version_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.goose_db_version_id_seq OWNED BY public.goose_db_version.id;


--
-- Name: meta_photo_data; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.meta_photo_data (
    photo_id uuid NOT NULL,
    model_info text,
    size_bytes integer NOT NULL,
    width_pixel integer NOT NULL,
    height_pixel integer NOT NULL,
    date_time timestamp with time zone,
    updated_at timestamp with time zone NOT NULL,
    geo_latitude double precision,
    geo_longitude double precision,
    CONSTRAINT meta_photo_data_model_info_check CHECK ((length(model_info) <= 512))
);


--
-- Name: people; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.people (
    id uuid NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    firstname text NOT NULL,
    surname text NOT NULL,
    patronymic text,
    CONSTRAINT people_firstname_check CHECK ((length(firstname) <= 1024)),
    CONSTRAINT people_patronymic_check CHECK ((length(patronymic) <= 1024)),
    CONSTRAINT people_surname_check CHECK ((length(surname) <= 1024))
);


--
-- Name: photo_groups; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.photo_groups (
    id uuid NOT NULL,
    main_photo_id uuid NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL
);


--
-- Name: photo_groups_photos; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.photo_groups_photos (
    photo_id uuid NOT NULL,
    group_id uuid NOT NULL
);


--
-- Name: photo_locations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.photo_locations (
    photo_id uuid NOT NULL,
    created_at timestamp with time zone NOT NULL,
    geo_latitude double precision NOT NULL,
    geo_longitude double precision NOT NULL,
    formatted_address text NOT NULL,
    street text NOT NULL,
    house_number text NOT NULL,
    suburb text NOT NULL,
    postcode text NOT NULL,
    state text NOT NULL,
    state_code text NOT NULL,
    state_district text NOT NULL,
    county text NOT NULL,
    country text NOT NULL,
    country_code text NOT NULL,
    city text NOT NULL,
    CONSTRAINT photo_locations_city_check CHECK ((length(city) <= 1024)),
    CONSTRAINT photo_locations_country_check CHECK ((length(country) <= 1024)),
    CONSTRAINT photo_locations_country_code_check CHECK ((length(country_code) <= 1024)),
    CONSTRAINT photo_locations_county_check CHECK ((length(county) <= 1024)),
    CONSTRAINT photo_locations_formatted_address_check CHECK ((length(formatted_address) <= 1024)),
    CONSTRAINT photo_locations_house_number_check CHECK ((length(house_number) <= 1024)),
    CONSTRAINT photo_locations_postcode_check CHECK ((length(postcode) <= 1024)),
    CONSTRAINT photo_locations_state_check CHECK ((length(state) <= 1024)),
    CONSTRAINT photo_locations_state_code_check CHECK ((length(state_code) <= 1024)),
    CONSTRAINT photo_locations_state_district_check CHECK ((length(state_district) <= 1024)),
    CONSTRAINT photo_locations_street_check CHECK ((length(street) <= 1024)),
    CONSTRAINT photo_locations_suburb_check CHECK ((length(suburb) <= 1024))
);


--
-- Name: photo_previews; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.photo_previews (
    id uuid NOT NULL,
    photo_id uuid NOT NULL,
    file_key text NOT NULL,
    size_pixel integer NOT NULL,
    width_pixel integer NOT NULL,
    height_pixel integer NOT NULL,
    original boolean DEFAULT false NOT NULL,
    CONSTRAINT photo_previews_file_key_check CHECK ((length(file_key) <= 1024))
);


--
-- Name: photo_processing; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.photo_processing (
    photo_id uuid NOT NULL,
    processed_at timestamp with time zone NOT NULL,
    type public.processing_type NOT NULL,
    success boolean NOT NULL
);


--
-- Name: photo_tags; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.photo_tags (
    id uuid NOT NULL,
    category_id uuid NOT NULL,
    photo_id uuid NOT NULL,
    name text NOT NULL,
    CONSTRAINT photo_tags_name_check CHECK ((length(name) <= 128))
);


--
-- Name: photo_upload_data; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.photo_upload_data (
    photo_id uuid NOT NULL,
    upload_at timestamp with time zone NOT NULL,
    paths text[] NOT NULL,
    client_info text NOT NULL,
    person_id uuid NOT NULL,
    CONSTRAINT photo_upload_data_client_info_check CHECK ((length(client_info) <= 256)),
    CONSTRAINT photo_upload_data_paths_check CHECK ((cardinality(paths) <= 2048))
);


--
-- Name: photo_vectors; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.photo_vectors (
    photo_id uuid NOT NULL,
    vector double precision[] NOT NULL,
    norm double precision NOT NULL
);


--
-- Name: photos; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.photos (
    id uuid NOT NULL,
    file_key text NOT NULL,
    hash text NOT NULL,
    photo_updated_at timestamp with time zone NOT NULL,
    extension public.photo_extension NOT NULL,
    status public.photo_status DEFAULT 'ACTIVE'::public.photo_status NOT NULL,
    error text,
    CONSTRAINT photos_file_key_check CHECK ((length(file_key) <= 1024)),
    CONSTRAINT photos_hash_check CHECK ((length(hash) <= 512))
);


--
-- Name: refresh_codes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.refresh_codes (
    id uuid NOT NULL,
    person_id uuid NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    status public.refresh_token_status NOT NULL
);


--
-- Name: rocket_locks; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.rocket_locks (
    key text NOT NULL,
    locked_until timestamp with time zone NOT NULL,
    CONSTRAINT rocket_locks_key_check CHECK ((length(key) <= 128))
);


--
-- Name: tag_categories; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tag_categories (
    id uuid NOT NULL,
    type text NOT NULL,
    color text NOT NULL,
    CONSTRAINT tag_categories_color_check CHECK ((length(color) <= 7)),
    CONSTRAINT tag_categories_type_check CHECK ((length(type) <= 64))
);


--
-- Name: goose_db_version id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.goose_db_version ALTER COLUMN id SET DEFAULT nextval('public.goose_db_version_id_seq'::regclass);


--
-- Name: api_tokens api_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.api_tokens
    ADD CONSTRAINT api_tokens_pkey PRIMARY KEY (id);


--
-- Name: api_tokens api_tokens_token_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.api_tokens
    ADD CONSTRAINT api_tokens_token_key UNIQUE (token);


--
-- Name: auth auth_email_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.auth
    ADD CONSTRAINT auth_email_key UNIQUE (email);


--
-- Name: auth auth_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.auth
    ADD CONSTRAINT auth_pkey PRIMARY KEY (person_id);


--
-- Name: codes codes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.codes
    ADD CONSTRAINT codes_pkey PRIMARY KEY (code);


--
-- Name: coefficients_similar_photos coefficients_similar_photos_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.coefficients_similar_photos
    ADD CONSTRAINT coefficients_similar_photos_pkey PRIMARY KEY (photo_id1, photo_id2);


--
-- Name: exif_photo_data exif_photo_data_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.exif_photo_data
    ADD CONSTRAINT exif_photo_data_pkey PRIMARY KEY (photo_id);


--
-- Name: goose_db_version goose_db_version_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.goose_db_version
    ADD CONSTRAINT goose_db_version_pkey PRIMARY KEY (id);


--
-- Name: meta_photo_data meta_photo_data_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.meta_photo_data
    ADD CONSTRAINT meta_photo_data_pkey PRIMARY KEY (photo_id);


--
-- Name: people people_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.people
    ADD CONSTRAINT people_pkey PRIMARY KEY (id);


--
-- Name: photo_groups photo_groups_main_photo_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.photo_groups
    ADD CONSTRAINT photo_groups_main_photo_id_key UNIQUE (main_photo_id);


--
-- Name: photo_groups_photos photo_groups_photos_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.photo_groups_photos
    ADD CONSTRAINT photo_groups_photos_pkey PRIMARY KEY (photo_id, group_id);


--
-- Name: photo_groups photo_groups_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.photo_groups
    ADD CONSTRAINT photo_groups_pkey PRIMARY KEY (id);


--
-- Name: photo_locations photo_locations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.photo_locations
    ADD CONSTRAINT photo_locations_pkey PRIMARY KEY (photo_id);


--
-- Name: photo_previews photo_previews_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.photo_previews
    ADD CONSTRAINT photo_previews_pkey PRIMARY KEY (id);


--
-- Name: photo_processing photo_processing_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.photo_processing
    ADD CONSTRAINT photo_processing_pkey PRIMARY KEY (photo_id, type);


--
-- Name: photo_tags photo_tags_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.photo_tags
    ADD CONSTRAINT photo_tags_pkey PRIMARY KEY (id);


--
-- Name: photo_upload_data photo_upload_data_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.photo_upload_data
    ADD CONSTRAINT photo_upload_data_pkey PRIMARY KEY (photo_id);


--
-- Name: photo_vectors photo_vectors_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.photo_vectors
    ADD CONSTRAINT photo_vectors_pkey PRIMARY KEY (photo_id);


--
-- Name: photos photos_file_key_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.photos
    ADD CONSTRAINT photos_file_key_key UNIQUE (file_key);


--
-- Name: photos photos_hash_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.photos
    ADD CONSTRAINT photos_hash_key UNIQUE (hash);


--
-- Name: photos photos_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.photos
    ADD CONSTRAINT photos_pkey PRIMARY KEY (id);


--
-- Name: refresh_codes refresh_codes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.refresh_codes
    ADD CONSTRAINT refresh_codes_pkey PRIMARY KEY (id);


--
-- Name: rocket_locks rocket_locks_key_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.rocket_locks
    ADD CONSTRAINT rocket_locks_key_key UNIQUE (key);


--
-- Name: tag_categories tag_categories_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tag_categories
    ADD CONSTRAINT tag_categories_pkey PRIMARY KEY (id);


--
-- Name: tag_categories tag_categories_type_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tag_categories
    ADD CONSTRAINT tag_categories_type_key UNIQUE (type);


--
-- Name: idx_api_tokens_person_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_api_tokens_person_id ON public.api_tokens USING btree (person_id);


--
-- Name: idx_coefficients_similar_photos_photo_id2; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_coefficients_similar_photos_photo_id2 ON public.coefficients_similar_photos USING btree (photo_id2);


--
-- Name: idx_photo_previews_file_key; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_photo_previews_file_key ON public.photo_previews USING btree (file_key);


--
-- Name: idx_photo_previews_photo_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_photo_previews_photo_id ON public.photo_previews USING btree (photo_id);


--
-- Name: idx_photos_groups_references_group_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_photos_groups_references_group_id ON public.photo_groups_photos USING btree (group_id);


--
-- Name: idx_tags_photo_id_name; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_tags_photo_id_name ON public.photo_tags USING btree (photo_id, name);


--
-- Name: api_tokens api_tokens_person_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.api_tokens
    ADD CONSTRAINT api_tokens_person_id_fkey FOREIGN KEY (person_id) REFERENCES public.people(id);


--
-- Name: auth auth_person_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.auth
    ADD CONSTRAINT auth_person_id_fkey FOREIGN KEY (person_id) REFERENCES public.people(id);


--
-- Name: codes codes_person_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.codes
    ADD CONSTRAINT codes_person_id_fkey FOREIGN KEY (person_id) REFERENCES public.people(id);


--
-- Name: coefficients_similar_photos coefficients_similar_photos_photo_id1_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.coefficients_similar_photos
    ADD CONSTRAINT coefficients_similar_photos_photo_id1_fkey FOREIGN KEY (photo_id1) REFERENCES public.photos(id);


--
-- Name: coefficients_similar_photos coefficients_similar_photos_photo_id2_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.coefficients_similar_photos
    ADD CONSTRAINT coefficients_similar_photos_photo_id2_fkey FOREIGN KEY (photo_id2) REFERENCES public.photos(id);


--
-- Name: exif_photo_data exif_photo_data_photo_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.exif_photo_data
    ADD CONSTRAINT exif_photo_data_photo_id_fkey FOREIGN KEY (photo_id) REFERENCES public.photos(id);


--
-- Name: meta_photo_data meta_photo_data_photo_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.meta_photo_data
    ADD CONSTRAINT meta_photo_data_photo_id_fkey FOREIGN KEY (photo_id) REFERENCES public.photos(id);


--
-- Name: photo_groups photo_groups_main_photo_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.photo_groups
    ADD CONSTRAINT photo_groups_main_photo_id_fkey FOREIGN KEY (main_photo_id) REFERENCES public.photos(id);


--
-- Name: photo_groups_photos photo_groups_photos_group_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.photo_groups_photos
    ADD CONSTRAINT photo_groups_photos_group_id_fkey FOREIGN KEY (group_id) REFERENCES public.photo_groups(id);


--
-- Name: photo_groups_photos photo_groups_photos_photo_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.photo_groups_photos
    ADD CONSTRAINT photo_groups_photos_photo_id_fkey FOREIGN KEY (photo_id) REFERENCES public.photos(id);


--
-- Name: photo_locations photo_locations_photo_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.photo_locations
    ADD CONSTRAINT photo_locations_photo_id_fkey FOREIGN KEY (photo_id) REFERENCES public.photos(id);


--
-- Name: photo_previews photo_previews_photo_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.photo_previews
    ADD CONSTRAINT photo_previews_photo_id_fkey FOREIGN KEY (photo_id) REFERENCES public.photos(id);


--
-- Name: photo_processing photo_processing_photo_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.photo_processing
    ADD CONSTRAINT photo_processing_photo_id_fkey FOREIGN KEY (photo_id) REFERENCES public.photos(id);


--
-- Name: photo_tags photo_tags_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.photo_tags
    ADD CONSTRAINT photo_tags_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.tag_categories(id);


--
-- Name: photo_tags photo_tags_photo_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.photo_tags
    ADD CONSTRAINT photo_tags_photo_id_fkey FOREIGN KEY (photo_id) REFERENCES public.photos(id);


--
-- Name: photo_upload_data photo_upload_data_person_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.photo_upload_data
    ADD CONSTRAINT photo_upload_data_person_id_fkey FOREIGN KEY (person_id) REFERENCES public.people(id);


--
-- Name: photo_upload_data photo_upload_data_photo_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.photo_upload_data
    ADD CONSTRAINT photo_upload_data_photo_id_fkey FOREIGN KEY (photo_id) REFERENCES public.photos(id);


--
-- Name: photo_vectors photo_vectors_photo_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.photo_vectors
    ADD CONSTRAINT photo_vectors_photo_id_fkey FOREIGN KEY (photo_id) REFERENCES public.photos(id);


--
-- Name: refresh_codes refresh_codes_person_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.refresh_codes
    ADD CONSTRAINT refresh_codes_person_id_fkey FOREIGN KEY (person_id) REFERENCES public.people(id);


--
-- PostgreSQL database dump complete
--

