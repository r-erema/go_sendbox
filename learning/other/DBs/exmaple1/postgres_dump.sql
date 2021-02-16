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

CREATE TABLE public.edges (
    edge_id integer NOT NULL,
    tail_vertex integer,
    head_vertex integer,
    label text,
    properties json
);


ALTER TABLE public.edges OWNER TO go;

CREATE TABLE public.vertices (
    vertex_id integer NOT NULL,
    properties json
);


ALTER TABLE public.vertices OWNER TO go;

INSERT INTO public.edges (edge_id, tail_vertex, head_vertex, label, properties) VALUES (1, 2, 1, 'WITHIN', NULL);
INSERT INTO public.edges (edge_id, tail_vertex, head_vertex, label, properties) VALUES (2, 3, 2, 'WITHIN', NULL);
INSERT INTO public.edges (edge_id, tail_vertex, head_vertex, label, properties) VALUES (3, 5, 4, 'WITHIN', NULL);
INSERT INTO public.edges (edge_id, tail_vertex, head_vertex, label, properties) VALUES (4, 6, 4, 'WITHIN', NULL);
INSERT INTO public.edges (edge_id, tail_vertex, head_vertex, label, properties) VALUES (5, 7, 5, 'WITHIN', NULL);
INSERT INTO public.edges (edge_id, tail_vertex, head_vertex, label, properties) VALUES (6, 8, 7, 'WITHIN', NULL);
INSERT INTO public.edges (edge_id, tail_vertex, head_vertex, label, properties) VALUES (7, 9, 5, 'WITHIN', NULL);
INSERT INTO public.edges (edge_id, tail_vertex, head_vertex, label, properties) VALUES (8, 10, 9, 'WITHIN', NULL);
INSERT INTO public.edges (edge_id, tail_vertex, head_vertex, label, properties) VALUES (9, 11, 10, 'WITHIN', NULL);
INSERT INTO public.edges (edge_id, tail_vertex, head_vertex, label, properties) VALUES (10, 12, 3, 'BORN_IN', NULL);
INSERT INTO public.edges (edge_id, tail_vertex, head_vertex, label, properties) VALUES (11, 12, 8, 'LIVES_IN', NULL);
INSERT INTO public.edges (edge_id, tail_vertex, head_vertex, label, properties) VALUES (12, 13, 11, 'BORN_IN', NULL);
INSERT INTO public.edges (edge_id, tail_vertex, head_vertex, label, properties) VALUES (13, 13, 8, 'LIVES_IN', NULL);

INSERT INTO public.vertices (vertex_id, properties) VALUES (1, '{"name":"North America","type":"continent"}');
INSERT INTO public.vertices (vertex_id, properties) VALUES (2, '{"name":"United States","type":"country"}');
INSERT INTO public.vertices (vertex_id, properties) VALUES (3, '{"name":"Idaho","type":"state"}');
INSERT INTO public.vertices (vertex_id, properties) VALUES (4, '{"name":"Europe","type":"continent"}');
INSERT INTO public.vertices (vertex_id, properties) VALUES (5, '{"name":"United Kingdom","type":"country"}');
INSERT INTO public.vertices (vertex_id, properties) VALUES (6, '{"name":"France","type":"country"}');
INSERT INTO public.vertices (vertex_id, properties) VALUES (7, '{"name":"England","type":"country"}');
INSERT INTO public.vertices (vertex_id, properties) VALUES (8, '{"name":"London","type":"city"}');
INSERT INTO public.vertices (vertex_id, properties) VALUES (9, '{"name":"Burgundy","type":"region", "name_fr":"Bourgogne", "name_en": "Burgundy"}');
INSERT INTO public.vertices (vertex_id, properties) VALUES (10, '{"name":"Côte d\\''Ivoire","type":"department"}');
INSERT INTO public.vertices (vertex_id, properties) VALUES (11, '{"name":"Beaune","type":"city"}');
INSERT INTO public.vertices (vertex_id, properties) VALUES (12, '{"name":"Lucy"}');
INSERT INTO public.vertices (vertex_id, properties) VALUES (13, '{"name":"Alain"}');

ALTER TABLE ONLY public.edges
    ADD CONSTRAINT edges_pkey PRIMARY KEY (edge_id);

ALTER TABLE ONLY public.vertices
    ADD CONSTRAINT vertices_pkey PRIMARY KEY (vertex_id);

CREATE INDEX edges_heads ON public.edges USING btree (head_vertex);

CREATE INDEX edges_tails ON public.edges USING btree (tail_vertex);

ALTER TABLE ONLY public.edges
    ADD CONSTRAINT edges_head_vertex_fkey FOREIGN KEY (head_vertex) REFERENCES public.vertices(vertex_id);

ALTER TABLE ONLY public.edges
    ADD CONSTRAINT edges_tail_vertex_fkey FOREIGN KEY (tail_vertex) REFERENCES public.vertices(vertex_id);
