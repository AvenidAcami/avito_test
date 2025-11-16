CREATE TABLE public.members (
    user_id text NOT NULL,
    username text NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    team_name text,
    pull_request_id text
);


CREATE TABLE public.pull_requests (
    pull_request_id text NOT NULL,
    pull_request_name text NOT NULL,
    author_id text NOT NULL,
    status text NOT NULL DEFAULT 'OPEN',
    "createdAt" TIMESTAMP WITH TIME ZONE,
    "mergedAt" TIMESTAMP WITH TIME ZONE
);


CREATE TABLE public.teams (
    team_name text NOT NULL
);


ALTER TABLE ONLY public.members
    ADD CONSTRAINT members_pkey PRIMARY KEY (user_id);

ALTER TABLE ONLY public.pull_requests
    ADD CONSTRAINT pull_requests_pkey PRIMARY KEY (pull_request_id);

ALTER TABLE ONLY public.teams
    ADD CONSTRAINT teams_pkey PRIMARY KEY (team_name);

ALTER TABLE ONLY public.members
    ADD CONSTRAINT team FOREIGN KEY (team_name) REFERENCES public.teams(team_name) NOT VALID;


ALTER TABLE ONLY public.members
    ADD CONSTRAINT pull_request_id FOREIGN KEY (pull_request_id) REFERENCES public.pull_requests(pull_request_id) NOT VALID;