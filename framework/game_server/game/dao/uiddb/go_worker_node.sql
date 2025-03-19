DROP TABLE IF EXISTS  go_worker_node;
CREATE TABLE go_worker_node (
  id          bigint(20) NOT NULL comment 'id',
  host_name   varchar(64) comment '主机ip',
  create_time datetime NULL comment '创建时间',
  PRIMARY KEY (id)) comment='工作节点';