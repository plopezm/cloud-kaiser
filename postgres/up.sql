DROP TABLE IF EXISTS tasks;
CREATE TABLE tasks (
  name                    VARCHAR NOT NULL,
  version                 VARCHAR NOT NULL,
  created_at              TIMESTAMP WITH TIME ZONE NOT NULL,
  script                  VARCHAR NOT NULL,
  PRIMARY KEY (name, version)
);

DROP TABLE IF EXISTS jobs;
CREATE TABLE jobs (
  name            VARCHAR NOT NULL,
  version         VARCHAR NOT NULL,
  created_at      TIMESTAMP WITH TIME ZONE NOT NULL,
  activation_type VARCHAR NOT NULL,
  duration        VARCHAR,
  PRIMARY KEY(name, version)
);


DROP TABLE IF EXISTS arguments;
CREATE TABLE arguments (
  name            VARCHAR NOT NULL,
  job_name        VARCHAR NOT NULL,
  job_version     VARCHAR NOT NULL,
  value           VARCHAR,
  PRIMARY KEY (name, job_name, job_version),
  FOREIGN KEY (job_name, job_version) REFERENCES jobs(name, version)
);

DROP TABLE IF EXISTS jobs_tasks;
CREATE TABLE jobs_tasks (
  job_name                VARCHAR NOT NULL,
  job_version             VARCHAR NOT NULL,
  task_name               VARCHAR NOT NULL,
  task_version            VARCHAR NOT NULL,
  on_success_name         VARCHAR,
  on_success_version      VARCHAR,
  on_failure_name         VARCHAR,
  on_failure_version      VARCHAR,
  PRIMARY KEY (job_name, job_version, task_name, task_version),
  FOREIGN KEY (job_name, job_version) REFERENCES jobs(name, version),
  FOREIGN KEY (task_name, task_version) REFERENCES tasks(name, version),
  FOREIGN KEY (job_name, job_version, on_success_name, on_success_version) REFERENCES jobs_tasks(job_name, job_version, task_name, task_version),
  FOREIGN KEY (job_name, job_version, on_failure_name, on_failure_version) REFERENCES jobs_tasks(job_name, job_version, task_name, task_version)
);


