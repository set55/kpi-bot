--软件研发项目进度达成率--
select account ,AVG(progress_standard)

from (
select a.account,c.name ,c.end,c.realEnd,TIMESTAMPDIFF(DAY,c.end,c.realEnd) as 与预期相差, 
case when TIMESTAMPDIFF(DAY,c.end,c.realEnd) < 0 then "1.2" 
	when TIMESTAMPDIFF(DAY,c.end,c.realEnd) = 0 then "1.0" 
	when TIMESTAMPDIFF(DAY,c.end,c.realEnd) < 5 then "0.5" 
else "0" end as progress_standard
from zt_user a 
inner join zt_team b on b.account = a.account 
inner join zt_project c on c.type in("project","sprint") and c.id = b.root and c.status = "closed" and c.acl = "private" 
where a.account in ("set.su","paul.gao","justin.lee","samy.gou","champion.fu","alan.tin","shiwen.tin") and c.realEnd between date_add(curdate(), interval - day(curdate()) + 1 day) and last_day(curdate())
order by a.account,c.realEnd desc
) tmp
group by account

--需求达成率--

select tmp.account,sum(tmp.get_score) 
from ( 
	select DISTINCT c.id,a.account,c.title,c.estimate, 
		case when c.estimate between 0.1 and 3 then 0.5
			when c.estimate between 4 and 7 then 1.5
			when c.estimate between 8 and 15 then 3
			when c.estimate between 16 and 31 then 4.5 
			when c.estimate between 32 and 71 then 6
			when c.estimate between 72 and 144 then 7.5 
			else 0 end as get_score 
	from zt_user a 
	inner join zt_task b on b.finishedDate between date_add(curdate(), interval - day(curdate()) + 1 day) and last_day(curdate()) and b.finishedBy=a.account and b.deleted="0" and b.parent=0 
	inner join zt_story c on c.id = b.story and c.stage not in ("waiting","planned","projected","developing") 
	where a.account in ("set.su","paul.gao","justin.lee","samy.gou","champion.fu","alan.tin","shiwen.tin") order by a.account desc 
) tmp 
group by account

--需求达成率-完成情况--
select  account , get_score,count(1)

from 
(
select DISTINCT c.id,a.account,c.title,c.estimate, 
		case when c.estimate between 0 and 3 then 1
			when c.estimate between 4 and 7 then 1.5
			when c.estimate between 8 and 15 then 2
			else 2.5 end as get_score 
	from zt_user a 
	inner join zt_task b on b.finishedDate between date_add(curdate(), interval - day(curdate()) + 1 day) and last_day(curdate()) and b.finishedBy=a.account and b.deleted="0" and b.parent=0 
	inner join zt_story c on c.id = b.story and c.stage not in ("waiting","planned","projected","developing") 
	where a.account in ("set.su","paul.gao","justin.lee","samy.gou","champion.fu","alan.tin","shiwen.tin") order by a.account desc 
	
) tmp
GROUP BY account , get_score


--项目版本bug关闭率情况--（非嵌入式）

select account ,AVG(build_fix_bug/(build_postponed_bugs+build_fix_bug)) 
from ( 
	select a.account,c.title,c.project,c.execution,c.builds,count(1) as build_fix_bug,
	(select count(1) from zt_bug where openedBuild = c.builds and deleted = "0" and status = "active" and assignedTo = a.account) as build_postponed_bugs 
	from zt_user a 
	inner join zt_testreport c on c.end between date_add(curdate(), interval - day(curdate()) + 1 day) and last_day(curdate())and c.deleted ="0" 
	inner join zt_bug b on b.deleted="0" and b.project = c.project and b.openedBuild = c.builds and b.resolvedBy = a.account and b.resolution in ( "fixed") 
	where a.account in ("set.su","paul.gao","justin.lee","samy.gou","champion.fu","alan.tin","shiwen.tin") 
	group by a.account ,c.title,c.project,c.execution,c.builds 
	) tmp 
GROUP BY account


--项目版本bug关闭率情况---（嵌入式或无测试报告的）

select account ,AVG(build_fix_bug/(build_postponed_bugs+build_fix_bug)) 
from ( 
	select a.account,b.project,c.name,count(1) as build_fix_bug,
	(select count(1) from zt_bug where project = b.project and deleted = "0" and status = "active" and assignedTo = a.account and assignedDate between date_add(curdate(), interval - day(curdate()) + 1 day) and last_day(curdate())) as build_postponed_bugs 
	from zt_user a 
	inner join zt_bug b on b.deleted="0" and b.resolvedBy = a.account and b.resolution in ( "fixed") and assignedDate between date_add(curdate(), interval - day(curdate()) + 1 day) and last_day(curdate())
	inner join zt_project c on c.id = b.project
	where a.account in ("shiwen.tin")
	group by a.account ,b.project,c.name
) tmp 
GROUP BY account



--工时预估达成比--
select a.account,AVG(b.consumed/b.estimate) 
from zt_user a 
inner join zt_task b on b.finishedDate between date_add(curdate(), interval - day(curdate()) + 1 day) and last_day(curdate())and b.finishedBy=a.account and b.deleted="0" and b.parent=0 
where a.account in ("set.su","paul.gao","justin.lee","samy.gou","champion.fu","alan.tin","shiwen.tin") 
group by a.account 
order by a.account desc



--版本发版次数--

select a.account,b.project,b.execution,count(1) as pub_times 
from zt_user a 
inner join zt_testreport c on c.end between date_add(curdate(), interval - day(curdate()) + 1 day) and last_day(curdate())and c.deleted ="0" 
inner join zt_build b on b.id = c.builds and b.builder = a.account 
where a.account in ("set.su","paul.gao","justin.lee","samy.gou","champion.fu","alan.tin","shiwen.tin") 
GROUP BY a.account,b.project,b.execution


--版本发版次数-----（嵌入式或无测试报告的）

select a.account,b.project,b.execution,count(1) as pub_times 
from zt_user a 
inner join zt_build b on b.builder = a.account and b.createdDate between date_add(curdate(), interval - day(curdate()) + 1 day) and last_day(curdate())
where a.account in ("shiwen.tin") 
GROUP BY a.account,b.project,b.execution



--测试软件项目进度达成率--

select account,AVG(progress)
from (
select 

a.account,c.name,b.title,c.end,b.end as real_end,TIMESTAMPDIFF(DAY,b.end,c.end) as 与预期相差,
case when TIMESTAMPDIFF(DAY,b.end,c.end) < 0 then "1.2"
	when TIMESTAMPDIFF(DAY,b.end,c.end) = 0 then "1.0"
	when TIMESTAMPDIFF(DAY,b.end,c.end) < 5 then "0.5"
	else "0" end as progress
from zt_user a
inner join zt_testreport b on b.createdBy = a.account and b.end between date_add(curdate(), interval - day(curdate()) + 1 day) and last_day(curdate()) and b.deleted ="0"
inner join zt_testtask c on c.id = b.tasks
where a.account in ("linyanhai")
) tmp



--测试软件项目有效bug率--
--1、测试报告结束时间是当月的
--2、bug未被删除，bug关联项目属于测试报告关联项目，bug关联版本是测试报告所属版本，bug是焰海打开的，bug解决状态是转需求，延期处理和已解决的，这些叫有效bug。
--3、版本内所有bug，为项目与测试报告相等，并且不是指派给黄卫旗

select account ,AVG(build_fix_bugs/project_build_bugs) from 
(

select tmp.account,tmp.project,tmp.title, count(tmp.id) as build_fix_bugs,sum(DISTINCT tmp.build_bugs) as project_build_bugs
from (

select a.account,c.title,c.project,c.builds,b.id,
	 (select count(1) from zt_bug where project = c.project and openedBuild = c.builds and deleted = "0" and assignedTo not in ("huangweiqi") and  openedBy = a.account ) as build_bugs 
from zt_user a 
inner join zt_testreport c on c.end between date_add(curdate(), interval - day(curdate()) + 1 day) and last_day(curdate()) and c.deleted ="0" 
left join zt_bug b on b.deleted="0" and b.project = c.project and b.openedBuild = c.builds and b.openedBy = a.account and b.resolution in ("tostory","postponed","willnotfix","fixed") 
where a.account in ("linyanhai") 
) tmp 
group by tmp.account,tmp.project,tmp.title

) tmp2
--bug转需求数--

select a.account,b.id,b.title 
from zt_user a 
inner join zt_testreport c on c.end between date_add(curdate(), interval - day(curdate()) + 1 day) and last_day(curdate())and c.deleted ="0" 
inner join zt_bug b on b.deleted="0" and b.project = c.project and b.openedBuild = c.builds and b.openedBy = a.account and ( b.resolution in ("tostory")  or b.assignedTo in ("shawn.wang" , "huangweiqi"))
where a.account in ("linyanhai")




--用例发现bug率--同build下，有多少是关联case的
---"duplicate"重复bug,"bydesign"设计如此,"notrepro"未复现--
select account ,AVG(build_case_bugs/project_build_bugs) 
from 
( 


	select tmp.account,tmp.project , count(tmp.id) as build_case_bugs,sum(DISTINCT tmp.build_bugs) as project_build_bugs
	from (

	select a.account,c.title,c.project,c.builds,b.id,
		 (select count(1) from zt_bug where project = c.project and openedBuild = c.builds and deleted = "0" and assignedTo not in ("huangweiqi") and  openedBy = a.account ) as build_bugs 
	from zt_user a 
	inner join zt_testreport c on c.end between date_add(curdate(), interval - day(curdate()) + 1 day) and last_day(curdate()) and c.deleted ="0" 
	left join zt_bug b on b.deleted="0" and  b.project = c.project and  b.case <> 0 and b.openedBuild = c.builds and b.openedBy = a.account and b.resolution in ("tostory","postponed","willnotfix","fixed") 
	where a.account in ("linyanhai") 
	) tmp 
	group by tmp.account,tmp.project

	
) tmp2





--项目软件项目进度达成率和预估承诺达成率--

select account ,AVG((rd_standard+test_standard)/2) 
from ( 
	select a.account,b.name,b.type,b.end as 研发项目预期结束,b.realEnd as 研发项目实际结束,TIMESTAMPDIFF(DAY,b.end,b.realEnd) as 研发项目与预期相差, 
	case when TIMESTAMPDIFF(DAY,b.end,b.realEnd) < 0 then "1.2" 
		when TIMESTAMPDIFF(DAY,b.end,b.realEnd) = 0 then "1.0"
		when TIMESTAMPDIFF(DAY,b.end,b.realEnd) < 5 then "0.5" 
		else "0" end as rd_standard, d.end as 测试项目预期结束,
	c.end 测试项目实际结束,TIMESTAMPDIFF(DAY,c.end,c.createdDate) as 测试项目与预期相差, 
	case when TIMESTAMPDIFF(DAY,d.end,c.end) < 0 then "1.2" 
		when TIMESTAMPDIFF(DAY,d.end,c.end) = 0 then "1.0" 
		when TIMESTAMPDIFF(DAY,d.end,c.end) < 5 then "0.5" 
		else "0" end as test_standard
	from zt_user a 
	inner join zt_project b on b.PM = a.account and b.deleted="0" and b.realEnd between date_add(curdate(), interval - day(curdate()) + 1 day) and last_day(curdate())
	inner join zt_testreport c on (b.type="sprint" and c.execution = b.id) or (c.execution=0 and b.type="project" and c.project=b.id) 
	inner join zt_testtask d on d.id = c.tasks 
	where a.account in ("shawn.wang") 
) tmp






--项目成果完成率,不需要关注执行，只需要看项目需求完成度，因为有执行一定有项目--

select account,AVG(devleoped_num/project_storys) 
from ( 
	select a.account,b.id,b.name,count(1) as devleoped_num,
	(select count(1) from zt_projectstory ztp 
		inner join zt_story zts on zts.id=ztp.story and zts.deleted="0" 
		where project = b.id) as project_storys 
	from zt_user a 
	inner join zt_project b on b.PM = a.account and b.deleted="0" and b.realEnd between date_add(curdate(), interval - day(curdate()) + 1 day) and last_day(curdate())and b.project = 0 
	inner join zt_projectstory d on d.project = b.id inner join zt_story c on c.id= d.story and c.deleted = "0" and c.stage not in ("waiting","planned","projected","developing") 
	where a.account in ("shawn.wang") 
	group by a.account,b.id 
) tmp
 
--完成情况--
select account,name,devleoped_num/project_storys

from (
	select a.account,b.id,b.name,count(1) as devleoped_num,
	(select count(1) from zt_projectstory ztp 
		inner join zt_story zts on zts.id=ztp.story and zts.deleted="0" 
		where project = b.id) as project_storys 
	from zt_user a 
	inner join zt_project b on b.PM = a.account and b.deleted="0" and b.realEnd between date_add(curdate(), interval - day(curdate()) + 1 day) and last_day(curdate())and b.project = 0 
	inner join zt_projectstory d on d.project = b.id inner join zt_story c on c.id= d.story and c.deleted = "0" and c.stage not in ("waiting","planned","projected","developing") 
	where a.account in ("shawn.wang") 
	group by a.account,b.id 
) tmp


--项目规划需求数--
with recursive cte as ( 
	select b.id,d.story 
	from zt_project b 
	inner join zt_projectstory d on d.project = b.id 
	where b.deleted="0" and b.realEnd between date_add(curdate(), interval - day(curdate()) + 1 day) and last_day(curdate())and b.project = 0 
	) 

select b.openedBy,b.stage ,count(1) 
from zt_story b 
inner join cte a on a.story = b.id 
where b.openedBy = "shawn.wang" and b.deleted = "0" and b.id= a.story 
group by b.openedBy,b.stage




--预估承诺完成率，只看项目，和最后一个测试报告--

select account,AVG(pm_standard) 
from ( 

select account,name,plan,plan_end,testreport_end,TIMESTAMPDIFF(DAY,plan_end,testreport_end) as 研发项目与预期相差, 
case when TIMESTAMPDIFF(DAY,plan_end,testreport_end) < 0 then "1.2" 
	when TIMESTAMPDIFF(DAY,plan_end,testreport_end) = 0 then "1.0" 
	when TIMESTAMPDIFF(DAY,plan_end,testreport_end) < 5 then "0.5" 
else "0" end as 研发项目进度达基数 
from ( 
	select a.account,b.name,REPLACE(c.plan,",","") as plan,d.end as plan_end,
	(select end from zt_testreport a where id = (select max(id) from zt_testreport where project = b.id and product = d.product ) ) as testreport_end 
	from zt_user a 
	inner join zt_project b on b.PM = a.account and b.deleted="0" and b.realEnd between date_add(curdate(), interval - day(curdate()) + 1 day) and last_day(curdate())
	inner join zt_projectproduct c on c.project = b.id 
	inner join zt_productplan d on d.id=REPLACE(c.plan,",","") 
	where a.account in ("shawn.wang") 
) tmp where testreport_end is not null

)tmp2 
group by account





	