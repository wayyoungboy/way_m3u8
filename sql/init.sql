-- for sqllite
CREATE TABLE work_info (
                           ID  INTEGER PRIMARY KEY,
                           name CHAR(250) NOT NULL   ,
                           url  TEXT NOT NULL   ,
                           save_dir TEXT NOT NULL   ,
                           state INTEGER NOT NULL ,
                           info TEXT NOT NULL   ,
                           create_time TIMESTAMP    DEFAULT CURRENT_TIMESTAMP,
                           update_time TIMESTAMP    DEFAULT CURRENT_TIMESTAMP
) ;

-- 常用sql
update work_info set state=1 where state=3;
update work_info set save_dir="/data1/media/av/" where save_dir="../../m3u8";
