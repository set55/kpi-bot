--zt_task-- 任务表
--zt_story-- 需求表
--zt_project-- 项目/衝刺表
--zt_user-- 用户表

--软件研发项目进度达成率--
select tmp.account, AVG(tmp.diff_expect) as avg_diff_expect,
case when AVG(tmp.diff_expect) <= -5 then 1.2
		when AVG(tmp.diff_expect) > -5 and AVG(tmp.diff_expect) <= 0 then 1.0
		when AVG(tmp.diff_expect) > 0 and AVG(tmp.diff_expect) <= 2 then 0.8
		when AVG(tmp.diff_expect) > 2 and AVG(tmp.diff_expect) <= 4 then 0.6
		when AVG(tmp.diff_expect) > 4 and AVG(tmp.diff_expect) <= 6 then 0.5
		else 0 end as avg_diff_expect_standard
from (
select a.account,c.name ,c.end,c.realEnd,TIMESTAMPDIFF(DAY,c.end,c.realEnd) as diff_expect
from zt_user a 
inner join zt_team b on b.account = a.account 
inner join zt_project c on c.type in("project","sprint") and c.id = b.root and c.status = "closed" and c.acl in ("open", "private")
where a.account in ("set.su","paul.gao","justin.lee","samy.gou","champion.fu","alan.tin","shiwen.tin") and c.realEnd between date_add(curdate(), interval - day(curdate()) + 1 day) and last_day(curdate())
order by a.account,c.realEnd desc
) tmp
group by account

select tmp.account, AVG(tmp.diff_expect) as avg_diff_expect,
case when AVG(tmp.diff_expect) <= -5 then 1.2
		when AVG(tmp.diff_expect) > -5 and AVG(tmp.diff_expect) <= 0 then 1.0
		when AVG(tmp.diff_expect) > 0 and AVG(tmp.diff_expect) <= 2 then 0.8
		when AVG(tmp.diff_expect) > 2 and AVG(tmp.diff_expect) <= 4 then 0.6
		when AVG(tmp.diff_expect) > 4 and AVG(tmp.diff_expect) <= 6 then 0.5
		else 0 end as avg_diff_expect_standard
from (
select a.account,c.name ,c.end,c.realEnd,TIMESTAMPDIFF(DAY,c.end,c.realEnd) as diff_expect
from zt_user a 
inner join zt_team b on b.account = a.account 
inner join zt_project c on c.type in("project","sprint") and c.id = b.root and c.status = "closed" and c.acl in ("open", "private")
where a.account in ("set.su","paul.gao","justin.lee","samy.gou","champion.fu","alan.tin","shiwen.tin") and c.realEnd between "2024-07-01 00:00:00" and "2024-07-31 23:59:59"
order by a.account,c.realEnd desc
) tmp
group by account


--软件研发项目进度达成率-完成情况--
select a.account, d.name as project_name ,c.name as project_sprint_name,c.end,c.realEnd,TIMESTAMPDIFF(DAY,c.end,c.realEnd) as diff_day
from zt_user a 
inner join zt_team b on b.account = a.account 
inner join zt_project c on c.type in("sprint") and c.id = b.root and c.status = "closed" and c.acl in ("open", "private")
left join zt_project d on d.id = c.project
where a.account in ("set.su","paul.gao","justin.lee","samy.gou","champion.fu","alan.tin","shiwen.tin") and c.realEnd between date_add(curdate(), interval - day(curdate()) + 1 day) and last_day(curdate())
order by a.account,c.realEnd desc

--需求达成率--

select tmp.account,sum(tmp.get_score) 
from ( 
	select DISTINCT c.id,a.account,c.title,c.estimate, 
        case when c.estimate < 4 then 1
            when c.estimate < 8 and c.estimate >= 4 then 1.5
            when c.estimate < 16 and c.estimate >= 8 then 2
            when c.estimate >= 16 then 2.5
			else 0 end as get_score 
	from zt_user a 
	inner join zt_task b on b.finishedDate between date_add(curdate(), interval - day(curdate()) + 1 day) and last_day(curdate()) and b.finishedBy=a.account and b.deleted="0" and b.parent=0 
	inner join zt_story c on c.id = b.story and c.stage not in ("waiting","planned","projected","developing") 
	where a.account in ("set.su","paul.gao","justin.lee","samy.gou","champion.fu","alan.tin","shiwen.tin") order by a.account desc 
) tmp 
group by account

--需求达成率-完成情况--
select DISTINCT c.id,a.account,c.title,c.estimate, c.stage,
    case when c.estimate < 4 then 1
            when c.estimate < 8 and c.estimate >= 4 then 1.5
            when c.estimate < 16 and c.estimate >= 8 then 2
            when c.estimate >= 16 then 2.5
			else 0 end as get_score  
from zt_user a 
inner join zt_task b on b.finishedDate between date_add(curdate(), interval - day(curdate()) + 1 day) and last_day(curdate()) and b.finishedBy=a.account and b.deleted="0" and b.parent=0 
inner join zt_story c on c.id = b.story and c.stage not in ("waiting","planned","projected","developing") 
where a.account in ("set.su","paul.gao","justin.lee","samy.gou","champion.fu","alan.tin","shiwen.tin") order by a.account desc 


--项目版本bug遗留率情况--（非嵌入式）

select account ,AVG(build_postponed_bugs/(build_postponed_bugs+build_fix_bug)) as bug_carry_over_rate,
case when AVG(build_postponed_bugs/(build_postponed_bugs+build_fix_bug)) = 0 then "1.2"
    when AVG(build_postponed_bugs/(build_postponed_bugs+build_fix_bug)) <= 0.1 then "1.0"
	when AVG(build_postponed_bugs/(build_postponed_bugs+build_fix_bug)) <= 0.2 then "0.9"
	when AVG(build_postponed_bugs/(build_postponed_bugs+build_fix_bug)) <= 0.3 then "0.8"
	when AVG(build_postponed_bugs/(build_postponed_bugs+build_fix_bug)) <= 0.4 then "0.6"
	when AVG(build_postponed_bugs/(build_postponed_bugs+build_fix_bug)) <= 0.5 then "0.5"
    else "0" end as bug_carry_over_rate_standard,
sum(tmp.build_postponed_bugs) as build_postponed_bugs,sum(build_fix_bug) as build_fix_bug
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
-----------
select tmp.account ,AVG(tmp.build_postponed_bugs/(tmp.build_postponed_bugs+build_fix_bug)) as bug_carry_over_rate,
case when AVG(tmp.build_postponed_bugs/(tmp.build_postponed_bugs+build_fix_bug)) = 0 then "1.2"
    when AVG(tmp.build_postponed_bugs/(tmp.build_postponed_bugs+build_fix_bug)) <= 0.1 then "1.0"
	when AVG(tmp.build_postponed_bugs/(tmp.build_postponed_bugs+build_fix_bug)) <= 0.2 then "0.9"
	when AVG(tmp.build_postponed_bugs/(tmp.build_postponed_bugs+build_fix_bug)) <= 0.3 then "0.8"
	when AVG(tmp.build_postponed_bugs/(tmp.build_postponed_bugs+build_fix_bug)) <= 0.4 then "0.6"
	when AVG(tmp.build_postponed_bugs/(tmp.build_postponed_bugs+build_fix_bug)) <= 0.5 then "0.5"
    else "0" end as bug_carry_over_rate_standard,
sum(tmp.build_postponed_bugs) as build_postponed_bugs,sum(build_fix_bug) as build_fix_bug
from (
	select a.account as account, c.title as test_report, b.id as bug_id, b.title as bug_title, b.resolution as bug_resolution, b.status as bug_status, count(1) as build_fix_bug,
	(select count(1) from zt_bug where openedBuild = c.builds and deleted = "0" and status = "active" and assignedTo = a.account) as build_postponed_bugs 
	from zt_testreport c
	left join  zt_bug b on b.deleted="0" and b.project = c.project and b.openedBuild = c.builds
	left join zt_user a on a.account = b.resolvedBy or (a.account = b.assignedTo and b.status="active")
	where a.account in ("set.su","paul.gao","justin.lee","samy.gou","champion.fu","alan.tin","shiwen.tin") and c.end between date_add(curdate(), interval - day(curdate()) + 1 day) and last_day(curdate()) and c.deleted ="0"
	group by a.account ,c.title,c.project,c.execution,c.builds
) tmp group by tmp.account



--項目版本bug遗留實際情況--（非嵌入式）
select a.account as account, c.title as test_report, b.id as bug_id, b.title as bug_title, b.resolution as bug_resolution, b.status as bug_status,
(select count(1) from zt_bug where openedBuild = c.builds and deleted = "0" and status = "active" and assignedTo = a.account) as build_postponed_bugs 
from zt_testreport c
inner join  zt_bug b on b.deleted="0" and b.project = c.project and b.openedBuild = c.builds
inner join zt_user a on a.account = b.resolvedBy or (a.account = b.assignedTo and b.status="active")
where a.account in ("set.su","paul.gao","justin.lee","samy.gou","champion.fu","alan.tin","shiwen.tin") and c.end between date_add(curdate(), interval - day(curdate()) + 1 day) and last_day(curdate()) and c.deleted ="0"

--项目版本bug遗留率情况---（嵌入式或无测试报告的）

select account ,AVG(build_postponed_bugs/(build_postponed_bugs+build_fix_bug)) as bug_carry_over_rate,
case when AVG(build_postponed_bugs/(build_postponed_bugs+build_fix_bug)) = 0 then "1.2"
    when AVG(build_postponed_bugs/(build_postponed_bugs+build_fix_bug)) <= 0.1 then "1.0"
    when AVG(build_postponed_bugs/(build_postponed_bugs+build_fix_bug)) <= 0.2 then "0.8"
    when AVG(build_postponed_bugs/(build_postponed_bugs+build_fix_bug)) <= 0.5 then "0.5"
    else "0" end as bug_carry_over_rate_standard
from ( 
	select a.account,b.project,c.name,count(1) as build_fix_bug,
	(select count(1) from zt_bug where project = b.project and deleted = "0" and status = "active" and assignedTo = a.account and assignedDate between date_add(curdate(), interval - day(curdate()) + 1 day) and last_day(curdate())) as build_postponed_bugs 
	from zt_user a 
	inner join zt_bug b on b.deleted="0" and b.resolvedBy = a.account and b.resolution in ( "fixed") and assignedDate between date_add(curdate(), interval - day(curdate()) + 1 day) and last_day(curdate())
	inner join zt_project c on c.id = b.project
	where a.account in ("set.su","paul.gao","justin.lee","samy.gou","champion.fu","alan.tin","shiwen.tin")
	group by a.account ,b.project,c.name
) tmp 
GROUP BY account

--項目版本bug遗留實際情況--（嵌入式或无测试报告的）
select a.account as account, c.name as project_name, b.id as bug_id, b.title as bug_title, b.resolution as bug_resolution, b.status as bug_status
from zt_bug b
left join zt_project c on c.id = b.project
left join zt_user a on a.account = b.resolvedBy or (a.account = b.assignedTo and b.status="active")
where a.account in ("set.su","paul.gao","justin.lee","samy.gou","champion.fu","alan.tin","shiwen.tin") and b.assignedDate between date_add(curdate(), interval - day(curdate()) + 1 day) and last_day(curdate()) and b.deleted ="0"




--工时预估达成比--
select a.account,AVG(b.consumed/b.estimate) as time_estimate_rate,
case when AVG(b.consumed/b.estimate) <= 0.8 then "1.2"
    when AVG(b.consumed/b.estimate) > 0.8 and AVG(b.consumed/b.estimate) <= 1 then "1.0"
	when AVG(b.consumed/b.estimate) > 1 and AVG(b.consumed/b.estimate) <= 1.2 then "0.8"
	when AVG(b.consumed/b.estimate) > 1.2 and AVG(b.consumed/b.estimate) <= 1.4 then "0.6"
	when AVG(b.consumed/b.estimate) > 1.4 and AVG(b.consumed/b.estimate) <= 1.6 then "0.4"
    else "0" end as time_estimate_level
from zt_user a 
inner join zt_task b on b.finishedDate between date_add(curdate(), interval - day(curdate()) + 1 day) and last_day(curdate())and b.finishedBy=a.account and b.deleted="0" and b.parent=0 
where a.account in ("set.su","paul.gao","justin.lee","samy.gou","champion.fu","alan.tin","shiwen.tin") 
group by a.account 
order by a.account desc


--版本发版次数平均发版次数--
select tmp2.account, AVG(tmp2.pub_times) as pub_times, 
case when AVG(tmp2.pub_times) > 0 and AVG(tmp2.pub_times) <= 1 then 1.2
			when AVG(tmp2.pub_times) > 1 and AVG(tmp2.pub_times) <= 2 then 1.0
			when AVG(tmp2.pub_times) > 2 and AVG(tmp2.pub_times) <= 3 then 0.8
			when AVG(tmp2.pub_times) > 3 and AVG(tmp2.pub_times) <= 4 then 0.6
			when AVG(tmp2.pub_times) > 4 and AVG(tmp2.pub_times) <= 5 then 0.5
			else 0 end as pub_times_level
from
(
    select tmp.account, tmp.project_type, tmp.project_name, tmp.pub_times, tmp.last_pub_time
    from
    (
        select a.account as account ,c.type as project_type, c.name as project_name,count(1) as pub_times, max(b.createdDate) as last_pub_time
        from zt_user a 
        inner join zt_build b on b.builder = a.account 
        left join zt_project c on b.execution = c.id
        where a.account in ("set.su","paul.gao","justin.lee","samy.gou","champion.fu","alan.tin","shiwen.tin") 
        GROUP BY a.account,b.project,b.execution
    ) tmp where tmp.last_pub_time between date_add(curdate(), interval - day(curdate()) + 1 day) and last_day(curdate())
) tmp2 group by tmp2.account

--版本发版次数详情--
select tmp.account, tmp.project_type, tmp.project_name, tmp.pub_times, tmp.last_pub_time
from
(
    select a.account as account ,c.type as project_type, c.name as project_name,count(1) as pub_times, max(b.createdDate) as last_pub_time
    from zt_user a 
    inner join zt_build b on b.builder = a.account 
    left join zt_project c on b.execution = c.id
    where a.account in ("set.su","paul.gao","justin.lee","samy.gou","champion.fu","alan.tin","shiwen.tin") 
    GROUP BY a.account,b.project,b.execution
) tmp where tmp.last_pub_time between date_add(curdate(), interval - day(curdate()) + 1 day) and last_day(curdate())


--版本发版次数-----（嵌入式或无测试报告的）不用

select a.account,b.project,b.execution,count(1) as pub_times,
case when count(1) < 3 then "1.2"
    when count(1) = 3 then "1.0"
    when count(1) > 3 and count(1) <=5 then "0.8"
    else "0.55" end as pub_times_level
from zt_user a 
inner join zt_build b on b.builder = a.account and b.createdDate between date_add(curdate(), interval - day(curdate()) + 1 day) and last_day(curdate())
where a.account in ("shiwen.tin") 
GROUP BY a.account,b.project,b.execution



--测试软件项目进度达成率--
select tmp.account,AVG(tmp.diff_days) as avg_diff_days,
case when AVG(tmp.diff_days) <= -3 then "1.2"
			when AVG(tmp.diff_days) > -3 and AVG(tmp.diff_days) <= 0 then "1.0"
			when AVG(tmp.diff_days) > 0 and AVG(tmp.diff_days) <= 2 then "0.8"
			when AVG(tmp.diff_days) > 2 and AVG(tmp.diff_days) <= 4 then "0.6"
			when AVG(tmp.diff_days) > 4 and AVG(tmp.diff_days) <= 5 then "0.5"
			else "0" end as progress
from (
select 
a.account,c.name,b.title,c.end,b.end as real_end,TIMESTAMPDIFF(DAY,c.end,b.end) as diff_days
from zt_user a
inner join zt_testreport b on b.createdBy = a.account and b.end between date_add(curdate(), interval - day(curdate()) + 1 day) and last_day(curdate()) and b.deleted ="0"
inner join zt_testtask c on c.id = b.tasks
where a.account in ("linyanhai")
) tmp



--测试软件项目有效bug率--
--1、测试报告结束时间是当月的
--2、bug未被删除，bug关联项目属于测试报告关联项目，bug关联版本是测试报告所属版本，bug是焰海打开的，bug解决状态是转需求，延期处理和已解决的，不予解决，这些叫有效bug。
--3、版本内所有bug，为项目与测试报告相等，并且不是指派给黄卫旗

select account ,AVG(build_fix_bugs/project_build_bugs) as validate_bug_rate,
case when AVG(build_fix_bugs/project_build_bugs) >= 0.95 then 1.2
	when AVG(build_fix_bugs/project_build_bugs) >= 0.9 and AVG(build_fix_bugs/project_build_bugs) < 0.95 then 1.0
	when AVG(build_fix_bugs/project_build_bugs) >= 0.8 and AVG(build_fix_bugs/project_build_bugs) < 0.9 then 0.8
	when AVG(build_fix_bugs/project_build_bugs) >= 0.7 and AVG(build_fix_bugs/project_build_bugs) < 0.8 then 0.5
	else 0 end as validate_bug_rate_standard
from 
(
select tmp.account,tmp.project,tmp.title, count(tmp.id) as build_fix_bugs,sum(DISTINCT tmp.build_bugs) as project_build_bugs
from (
select a.account,c.title,c.project,c.builds,b.id,
	 (select count(1) from zt_bug where project = c.project and openedBuild = c.builds and deleted = "0" and assignedTo not in ("huangweiqi") and  openedBy = a.account ) as build_bugs 
from zt_user a 
inner join zt_testreport c on c.end between date_add(curdate(), interval - day(curdate()) + 1 day) and last_day(curdate()) and c.deleted ="0" 
left join zt_bug b on b.deleted="0" and b.project = c.project and b.openedBuild = c.builds and b.openedBy = a.account and b.resolution in ("tostory","postponed","willnotfix","fixed") 
where a.account in ("linyanhai") and b.assignedTo not in ("huangweiqi")
) tmp 
group by tmp.account,tmp.project,tmp.title
) tmp2


--bug转需求数--
select tmp.account, count(1) as tostory_num
from (
select a.account,b.id,b.title
from zt_user a 
inner join zt_testreport c on c.end between date_add(curdate(), interval - day(curdate()) + 1 day) and last_day(curdate())and c.deleted ="0" 
inner join zt_bug b on b.deleted="0" and b.project = c.project and b.openedBuild = c.builds and b.openedBy = a.account and ( b.resolution in ("tostory")  or b.assignedTo in ("shawn.wang" , "huangweiqi"))
where a.account in ("linyanhai")
) tmp group by tmp.account




--用例发现bug率--同build(版本)下，有多少是关联case的
---"duplicate"重复bug,"bydesign"设计如此,"notrepro"未复现--
select account ,AVG(build_case_bugs/project_build_bugs)  as case_bug_rate,
case when AVG(build_case_bugs/project_build_bugs) >= 0.95 then 1.2
			when AVG(build_case_bugs/project_build_bugs) >= 0.9 and AVG(build_case_bugs/project_build_bugs) < 0.95 then 1.0
			when AVG(build_case_bugs/project_build_bugs) >= 0.8 and AVG(build_case_bugs/project_build_bugs) < 0.9 then 0.8
			when AVG(build_case_bugs/project_build_bugs) >= 0.7 and AVG(build_case_bugs/project_build_bugs) < 0.8 then 0.5
			else 0 end as case_bug_rate_standard
from 
( 
	select tmp.account,tmp.project , count(tmp.id) as build_case_bugs,sum(DISTINCT tmp.build_bugs) as project_build_bugs
	from (
	select a.account,c.title,c.project,c.builds,b.id,
		 (select count(1) from zt_bug where project = c.project and openedBuild = c.builds and deleted = "0" and assignedTo not in ("huangweiqi") and  openedBy = a.account ) as build_bugs 
	from zt_user a 
	inner join zt_testreport c on c.end between date_add(curdate(), interval - day(curdate()) + 1 day) and last_day(curdate()) and c.deleted ="0" 
	left join zt_bug b on b.deleted="0" and  b.project = c.project and  b.case <> 0 and b.openedBuild = c.builds and b.openedBy = a.account and b.resolution in ("tostory","postponed","willnotfix","fixed") 
	where a.account in ("linyanhai")  and b.assignedTo not in ("huangweiqi")
	) tmp 
	group by tmp.account,tmp.project
) tmp2





--项目软件项目进度达成率--

select account ,((AVG(project_diff) + AVG(test_diff))/2) as avg_diff_days,
case when ((AVG(project_diff) + AVG(test_diff))/2) <= -3 then 1.2
			when ((AVG(project_diff) + AVG(test_diff))/2) > -3 and ((AVG(project_diff) + AVG(test_diff))/2) <= 0 then 1.0
			when ((AVG(project_diff) + AVG(test_diff))/2) > 0 and ((AVG(project_diff) + AVG(test_diff))/2) <= 2 then 0.8
			when ((AVG(project_diff) + AVG(test_diff))/2) > 2 and ((AVG(project_diff) + AVG(test_diff))/2) <= 4 then 0.6
			when ((AVG(project_diff) + AVG(test_diff))/2) > 4 and ((AVG(project_diff) + AVG(test_diff))/2) <= 5 then 0.5
			else 0 end as progress_standard
from ( 
	select a.account,b.name,b.type,b.end as project_end,b.realEnd as project_real_end,TIMESTAMPDIFF(DAY,b.end,b.realEnd) as project_diff,
	c.end test_end,TIMESTAMPDIFF(DAY,c.end,c.createdDate) as test_diff
	from zt_user a 
	inner join zt_project b on b.PM = a.account and b.deleted="0" and b.realEnd between date_add(curdate(), interval - day(curdate()) + 1 day) and last_day(curdate())
	inner join zt_testreport c on (b.type="sprint" and c.execution = b.id) or (c.execution=0 and b.type="project" and c.project=b.id) 
	inner join zt_testtask d on d.id = c.tasks 
	where a.account in ("shawn.wang") 
) tmp






--项目成果完成率,不需要关注执行，只需要看项目需求完成度，因为有执行一定有项目--
select account,AVG(devleoped_num/project_storys) as complete_rate,
case when AVG(devleoped_num/project_storys) >= 0.95 then 1.2
			when AVG(devleoped_num/project_storys) >= 0.9 and AVG(devleoped_num/project_storys) < 0.95 then 1.0
			when AVG(devleoped_num/project_storys) >= 0.8 and AVG(devleoped_num/project_storys) < 0.9 then 0.8
			when AVG(devleoped_num/project_storys) >= 0.7 and AVG(devleoped_num/project_storys) < 0.8 then 0.5
			else 0 end as complete_rate_standard
from ( 
	select a.account,b.id,b.name,count(1) as devleoped_num,
	(select count(1) from zt_projectstory ztp 
		inner join zt_story zts on zts.id=ztp.story and zts.deleted="0" 
		where project = b.id) as project_storys 
	from zt_user a 
	inner join zt_project b on b.PM = a.account and b.deleted="0" and b.realEnd between date_add(curdate(), interval - day(curdate()) + 1 day) and last_day(curdate())and b.project = 0 
	inner join zt_projectstory d on d.project = b.id
    inner join zt_story c on c.id= d.story and c.deleted = "0" and c.stage not in ("waiting","planned","projected","developing") 
	where a.account in ("shawn.wang") 
	group by a.account,b.id 
) tmp
 
--项目成果完成率,完成情况--
select account,name as project_name,devleoped_num/project_storys as complete_rate
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
where b.openedBy in ("shawn.wang") and b.deleted = "0"
group by b.openedBy,b.stage




--预估承诺完成率，只看项目，和最后一个测试报告--

select tmp2.account,AVG(tmp2.diff_days) as diff_days
from ( 
	select account,name,plan,plan_end,testreport_end,TIMESTAMPDIFF(DAY,plan_end,testreport_end) as diff_days 
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





	