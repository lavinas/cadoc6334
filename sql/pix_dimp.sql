drop table if exists pix_dimp_nsu;

create table pix_dimp_nsu (
    nsu_source varchar(50) primary key,
    nsu_target varchar(20),
    auth_target varchar(6)
);

insert into pix_dimp_nsu (nsu_source, nsu_target, auth_target)
select nsu as nsu_source,
       concat('P', substr(nsu, 5, 6), substr(nsu, 15, 9), RIGHT(nsu, 4)) as nsu_target,
       right(nsu, 6) as auth_target
  from pix_dimp
where duplicated = false;

alter table pix_dimp add column auth_target varchar(6);
alter table pix_dimp add column nsu_target varchar(20); 

update pix_dimp pd
   set auth_target = right(nsu, 6),
       nsu_target = concat('P', substr(nsu, 5, 6), substr(nsu, 15, 9), RIGHT(nsu, 4));

select * from pix_dimp;