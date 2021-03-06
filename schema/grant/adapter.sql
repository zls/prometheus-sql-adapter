GRANT ALL ON metric_labels TO :role_name;
GRANT ALL ON metric_samples TO :role_name;
GRANT ALL ON metrics TO :role_name;

GRANT ALL ON agg_instance_load TO :role_name;
GRANT ALL ON agg_instance_load_long TO :role_name;
GRANT ALL ON agg_instance_pods TO :role_name;

GRANT ALL ON agg_container_cpu TO :role_name;
GRANT ALL ON agg_container_mem TO :role_name;

GRANT ALL ON agg_grafana_alert TO :role_name;
GRANT ALL ON agg_grafana_alert_long TO :role_name;

GRANT ALL ON cat_container TO :role_name;
GRANT ALL ON cat_instance TO :role_name;
GRANT ALL ON cat_name TO :role_name;