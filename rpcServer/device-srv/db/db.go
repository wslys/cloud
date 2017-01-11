package db

var (
	GatewaySchema = `CREATE TABLE IF NOT EXISTS gateway (
						id int(11) NOT NULL AUTO_INCREMENT,
						object_id varchar(36) DEFAULT NULL,
						name varchar(60) DEFAULT NULL,
						mac varchar(18) DEFAULT NULL,
						status int(11) DEFAULT NULL,
						last_time int(11) DEFAULT NULL,
						site varchar(255) DEFAULT NULL,
						create_at int(11) DEFAULT NULL,
						update_at int(11) DEFAULT NULL,
						PRIMARY KEY (id),
						UNIQUE KEY mac (mac) USING BTREE
					) ENGINE=InnoDB AUTO_INCREMENT=81 DEFAULT CHARSET=latin1;`

	BeaconSchema = `CREATE TABLE IF NOT EXISTS beacon (
					  id int(11) NOT NULL AUTO_INCREMENT,
					  object_id varchar(24) NOT NULL,
					  user_id varchar(100) NOT NULL,
					  mac varchar(18) NOT NULL,
					  status int(11) DEFAULT NULL,
					  password varchar(100) DEFAULT NULL,
					  change_password varchar(100) DEFAULT NULL,
					  type int(11) NOT NULL DEFAULT '0',
					  current_version int(11) DEFAULT '0',
					  last_setting_version int(11) DEFAULT '0',
					  apply_status int(11) DEFAULT '0',
					  create_at int(11) NOT NULL,
					  update_at int(11) NOT NULL,
					  PRIMARY KEY (id)
					) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=latin1;`

	BeaconSetting = `CREATE TABLE IF NOT EXISTS beacon_setting (
					  id int(11) NOT NULL AUTO_INCREMENT,
					  object_id varchar(24) NOT NULL,
					  mac varchar(18) NOT NULL,
					  version int(11) NOT NULL,
					  apply_status int(11) DEFAULT NULL,
					  setting json NOT NULL,
					  create_at int(11) NOT NULL,
					  apply_at int(11) DEFAULT NULL,
					  PRIMARY KEY (id)
					) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=latin1;`
)
