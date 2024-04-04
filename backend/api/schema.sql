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
-- Name: photo_processing_status; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.photo_processing_status AS ENUM (
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
-- Name: photo_status; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.photo_status AS ENUM (
    'NEW_PHOTO',
    'NOT_VALID'
);


SET default_table_access_method = heap;

--
-- Name: coeffs_similar_photos; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.coeffs_similar_photos (
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
-- Name: locations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.locations (
    photo_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    geo_latitude double precision,
    geo_longitude double precision,
    formatted_address text,
    street text,
    house_number text,
    suburb text,
    postcode text,
    state text,
    state_code text,
    state_district text,
    county text,
    country text,
    country_code text,
    city text
);


--
-- Name: photo_groups; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.photo_groups (
    id uuid NOT NULL,
    main_photo_id uuid NOT NULL,
    update_at timestamp without time zone NOT NULL
);


--
-- Name: photo_groups_photos; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.photo_groups_photos (
    photo_id uuid NOT NULL,
    group_id uuid NOT NULL
);


--
-- Name: photo_metadata; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.photo_metadata (
    photo_id uuid NOT NULL,
    model_info text,
    size_bytes integer NOT NULL,
    width_pixel integer NOT NULL,
    height_pixel integer NOT NULL,
    date_time timestamp without time zone,
    update_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    geo_latitude double precision,
    geo_longitude double precision
);


--
-- Name: photo_previews; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.photo_previews (
    id uuid NOT NULL,
    photo_id uuid NOT NULL,
    file_name character varying(2048) NOT NULL,
    width_pixel integer NOT NULL,
    height_pixel integer NOT NULL,
    size_pixel integer NOT NULL
);


--
-- Name: photo_processing_statuses; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.photo_processing_statuses (
    photo_id uuid NOT NULL,
    processed_at timestamp without time zone NOT NULL,
    status public.photo_processing_status NOT NULL,
    success boolean DEFAULT false
);


--
-- Name: photo_upload_data; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.photo_upload_data (
    upload_at timestamp without time zone NOT NULL,
    photo_id uuid NOT NULL,
    paths text[] NOT NULL,
    client_id character varying(256) NOT NULL,
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
    file_name character varying(2048) NOT NULL,
    hash character varying(512) NOT NULL,
    update_at timestamp without time zone NOT NULL,
    extension character varying(8) NOT NULL,
    status public.photo_status DEFAULT 'NEW_PHOTO'::public.photo_status,
    error text
);


--
-- Name: rocket_locks; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.rocket_locks (
    key text NOT NULL,
    locked_until timestamp with time zone NOT NULL
);


--
-- Name: tags; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tags (
    id uuid NOT NULL,
    category_id uuid NOT NULL,
    photo_id uuid NOT NULL,
    name character varying(128) NOT NULL
);


--
-- Name: tags_category; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tags_category (
    id uuid NOT NULL,
    type character varying(128) NOT NULL,
    color character varying(7) NOT NULL
);


--
-- Name: goose_db_version id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.goose_db_version ALTER COLUMN id SET DEFAULT nextval('public.goose_db_version_id_seq'::regclass);


--
-- Name: coeffs_similar_photos coeffs_similar_photos_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.coeffs_similar_photos
    ADD CONSTRAINT coeffs_similar_photos_pkey PRIMARY KEY (photo_id1, photo_id2);


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
-- Name: locations locations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.locations
    ADD CONSTRAINT locations_pkey PRIMARY KEY (photo_id);


--
-- Name: photo_groups_photos photo_groups_photos_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.photo_groups_photos
    ADD CONSTRAINT photo_groups_photos_pkey PRIMARY KEY (photo_id);


--
-- Name: photo_groups photo_groups_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.photo_groups
    ADD CONSTRAINT photo_groups_pkey PRIMARY KEY (id);


--
-- Name: photo_metadata photo_metadata_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.photo_metadata
    ADD CONSTRAINT photo_metadata_pkey PRIMARY KEY (photo_id);


--
-- Name: photo_previews photo_previews_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.photo_previews
    ADD CONSTRAINT photo_previews_pkey PRIMARY KEY (id);


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
-- Name: photos photos_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.photos
    ADD CONSTRAINT photos_pkey PRIMARY KEY (id);


--
-- Name: rocket_locks rocket_locks_key_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.rocket_locks
    ADD CONSTRAINT rocket_locks_key_key UNIQUE (key);


--
-- Name: tags_category tags_category_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tags_category
    ADD CONSTRAINT tags_category_pkey PRIMARY KEY (id);


--
-- Name: tags tags_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tags
    ADD CONSTRAINT tags_pkey PRIMARY KEY (id);


--
-- Name: idx_coeffs_similar_photos_photo_id2; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_coeffs_similar_photos_photo_id2 ON public.coeffs_similar_photos USING btree (photo_id2);


--
-- Name: idx_photo_groups_photos_group_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_photo_groups_photos_group_id ON public.photo_groups_photos USING btree (group_id);


--
-- Name: idx_photo_previews_photo_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_photo_previews_photo_id ON public.photo_previews USING btree (photo_id);


--
-- Name: idx_photo_processing_statuses_photo_id_status; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_photo_processing_statuses_photo_id_status ON public.photo_processing_statuses USING btree (photo_id, status);


--
-- Name: idx_photos_file_name; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_photos_file_name ON public.photos USING btree (file_name);


--
-- Name: idx_photos_hash; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_photos_hash ON public.photos USING btree (hash);


--
-- Name: idx_tags_category_type; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_tags_category_type ON public.tags_category USING btree (type);


--
-- Name: idx_tags_photo_id_name; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_tags_photo_id_name ON public.tags USING btree (photo_id, name);


--
-- Name: coeffs_similar_photos coeffs_similar_photos_photo_id1_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.coeffs_similar_photos
    ADD CONSTRAINT coeffs_similar_photos_photo_id1_fkey FOREIGN KEY (photo_id1) REFERENCES public.photos(id);


--
-- Name: coeffs_similar_photos coeffs_similar_photos_photo_id2_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.coeffs_similar_photos
    ADD CONSTRAINT coeffs_similar_photos_photo_id2_fkey FOREIGN KEY (photo_id2) REFERENCES public.photos(id);


--
-- Name: exif_photo_data fk_exif_photo_data_photo_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.exif_photo_data
    ADD CONSTRAINT fk_exif_photo_data_photo_id FOREIGN KEY (photo_id) REFERENCES public.photos(id) ON DELETE CASCADE;


--
-- Name: locations fk_locations_photo_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.locations
    ADD CONSTRAINT fk_locations_photo_id FOREIGN KEY (photo_id) REFERENCES public.photos(id) ON DELETE CASCADE;


--
-- Name: photo_metadata fk_photo_metadata_photo_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.photo_metadata
    ADD CONSTRAINT fk_photo_metadata_photo_id FOREIGN KEY (photo_id) REFERENCES public.photos(id) ON DELETE CASCADE;


--
-- Name: photo_processing_statuses fk_photo_processing_statuses_photo_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.photo_processing_statuses
    ADD CONSTRAINT fk_photo_processing_statuses_photo_id FOREIGN KEY (photo_id) REFERENCES public.photos(id);


--
-- Name: photo_upload_data fk_photo_upload_data_photo_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.photo_upload_data
    ADD CONSTRAINT fk_photo_upload_data_photo_id FOREIGN KEY (photo_id) REFERENCES public.photos(id) ON DELETE CASCADE;


--
-- Name: photo_vectors fk_photo_vector_photo_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.photo_vectors
    ADD CONSTRAINT fk_photo_vector_photo_id FOREIGN KEY (photo_id) REFERENCES public.photos(id);


--
-- Name: tags fk_tag_category_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tags
    ADD CONSTRAINT fk_tag_category_id FOREIGN KEY (category_id) REFERENCES public.tags_category(id) ON DELETE CASCADE;


--
-- Name: tags fk_tag_photo_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tags
    ADD CONSTRAINT fk_tag_photo_id FOREIGN KEY (photo_id) REFERENCES public.photos(id) ON DELETE CASCADE;


--
-- Name: locations locations_photo_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.locations
    ADD CONSTRAINT locations_photo_id_fkey FOREIGN KEY (photo_id) REFERENCES public.photos(id);


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
-- Name: photo_previews photo_previews_photo_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.photo_previews
    ADD CONSTRAINT photo_previews_photo_id_fkey FOREIGN KEY (photo_id) REFERENCES public.photos(id);


--
-- PostgreSQL database dump complete
--

