

CREATE TABLE Node (
  id bigint NOT NULL PRIMARY KEY,
  name char NOT NULL,
  kme_name char,
  coord_lat bigint NOT NULL,
  coord_long bigint NOT NULL,
  type char NOT NULL
);


CREATE TABLE User (
  id bigint NOT NULL PRIMARY KEY,
  email char NOT NULL,
  password char NOT NULL,
  role char NOT NULL,
  name char,
  affiliation char
);


CREATE TABLE Device (
  id bigint NOT NULL PRIMARY KEY,
  model_id bigint NOT NULL,
  node_id bigint
);


CREATE TABLE App (
  id bigint NOT NULL PRIMARY KEY,
  name char NOT NULL,
  cert_pubkey blob
);


CREATE TABLE DeviceModel (
  id bigint NOT NULL PRIMARY KEY,
  vendor char NOT NULL,
  model char NOT NULL,
  max_keyrate bigint NOT NULL
);


CREATE TABLE LogItem (
  id bigint NOT NULL PRIMARY KEY,
  message char NOT NULL,
  timestamp datetime NOT NULL,
  type char
);


CREATE TABLE DeviceToOtherEnd (
  id bigint NOT NULL PRIMARY KEY,
  device_id bigint,
  other_end_id bigint
);


CREATE TABLE AppToNode (
  id bigint NOT NULL PRIMARY KEY,
  app_id bigint NOT NULL,
  node_id bigint NOT NULL,
  keyrate blob
);


CREATE TABLE DeviceToLog (
  id bigint NOT NULL PRIMARY KEY,
  device_id bigint,
  log_id bigint
);


CREATE TABLE AppToLog (
  id bigint NOT NULL PRIMARY KEY,
  app_id bigint,
  log_id bigint
);


ALTER TABLE Device ADD CONSTRAINT Device_model_id_fk FOREIGN KEY (model_id) REFERENCES DeviceModel (id);
ALTER TABLE DeviceToOtherEnd ADD CONSTRAINT DeviceToOtherEnd_device_id_fk FOREIGN KEY (device_id) REFERENCES Device (id);
ALTER TABLE DeviceToOtherEnd ADD CONSTRAINT DeviceToOtherEnd_other_end_id_fk FOREIGN KEY (other_end_id) REFERENCES Device (id);
ALTER TABLE Device ADD CONSTRAINT Device_node_id_fk FOREIGN KEY (node_id) REFERENCES Node (id);
ALTER TABLE AppToNode ADD CONSTRAINT AppToNode_app_id_fk FOREIGN KEY (app_id) REFERENCES App (id);
ALTER TABLE AppToNode ADD CONSTRAINT AppToNode_node_id_fk FOREIGN KEY (node_id) REFERENCES Node (id);
ALTER TABLE DeviceToLog ADD CONSTRAINT DeviceToLog_device_id_fk FOREIGN KEY (device_id) REFERENCES Device (id);
ALTER TABLE DeviceToLog ADD CONSTRAINT DeviceToLog_log_id_fk FOREIGN KEY (log_id) REFERENCES LogItem (id);
ALTER TABLE AppToLog ADD CONSTRAINT AppToLog_log_id_fk FOREIGN KEY (log_id) REFERENCES LogItem (id);
ALTER TABLE AppToLog ADD CONSTRAINT AppToLog_app_id_fk FOREIGN KEY (app_id) REFERENCES App (id);
