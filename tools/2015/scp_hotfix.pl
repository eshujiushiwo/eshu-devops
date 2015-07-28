#!/usr/bin/perl -w
#================================================================ 
#  (C) 2013-2014 Dena Group Holding Limited.
# 
#  This program is used to operate the process of  
#  login & web server.
#  Please check the options before using.
#
#  Authors:   
#    Edison Chow <zhou.liyang@dena.jp>
#================================================================
use warnings;
use strict;
use Net::OpenSSH;
use Getopt::Long;
my $num;
my %opt;
GetOptions(\%opt,
		'h|help',
	  'n|num=i'
	) or &print_usage();
if (!scalar(%opt) ) 
{
		&print_usage();
}
$opt{'h'} and &print_usage();
$opt{'n'} and $num=$opt{'n'};
sub print_usage()
{
printf <<EOF;
#================================================================ 
#  (C) 2013-2014 Dena Group Holding Limited.
# 
#  This program is used to operate the process of  
#  login & web server.
#  Please check the options before using.
#
#  Authors:   
#    Edison Chow <zhou.liyang\@dena.jp>
#================================================================
=================================================================
   -h,--help           Print Help Info. 
   -n,--num            The num of  Staticdata.

Sample :
   shell > ./nba.pl -n 100 
=================================================================
EOF
exit;
}
#================================================================
#		Function	total
#================================================================
sub total() 
{
my $flag1;
my $flag2;
	print "将获取以下版本号的Hotfix到正服\n";
	print "$num\n";
print "\n确定吗? (yes/no): ";
my $ret = <STDIN>; chomp($ret);
die "abort\n" if ($ret ne 'yes');


my $ssh1 = Net::OpenSSH->new("nba2015_login1"); 
$ssh1->error and warn "can not connect to nba2015_login1" . $ssh1->error;
      my $ssh2 = Net::OpenSSH->new("nba2015_login2");
			      $ssh2->error and warn "can not connect to nba2015_login2" . $ssh2->error;

$ssh1->system("rsync -avz root\@10.96.36.52:/www/doc/1/ios/HotfixConfig/HotfixConfig_$num.unity3d /www/doc/1/ios/HotfixConfig/");
$ssh2->system("rsync -avz root\@10.96.36.52:/www/doc/1/ios/HotfixConfig/HotfixConfig_$num.unity3d /www/doc/1/ios/HotfixConfig/");
 $ssh1->system("rsync -avz root\@10.96.36.52:/www/doc/1/android/HotfixConfig/HotfixConfig_$num.unity3d /www/doc/1/android/HotfixConfig/");
 $ssh2->system("rsync -avz root\@10.96.36.52:/www/doc/1/android/HotfixConfig/HotfixConfig_$num.unity3d /www/doc/1/android/HotfixConfig/");
 $ssh1->system("rsync -avz root\@10.96.36.52:/www/doc/1/wp8/HotfixConfig/HotfixConfig_$num.unity3d /www/doc/1/wp8/HotfixConfig/");
 $ssh2->system("rsync -avz root\@10.96.36.52:/www/doc/1/wp8/HotfixConfig/HotfixConfig_$num.unity3d /www/doc/1/wp8/HotfixConfig/");
				
				
my $a1;
my $a2;
my $a3;
my $a4;
my $a5;
my $a6;

$a1=$ssh1->capture("cd /www/doc/1/ios/HotfixConfig/ && ls -l|grep -v grep|grep HotfixConfig_$num |wc -l ");
$a2=$ssh1->capture("cd /www/doc/1/android/HotfixConfig/ && ls -l|grep -v grep|grep HotfixConfig_$num |wc -l ");
$a3=$ssh1->capture("cd /www/doc/1/wp8/HotfixConfig/  && ls -l |grep -v grep |grep HotfixConfig_$num |wc -l ");

$a4=$ssh2->capture("cd /www/doc/1/ios/HotfixConfig/ && ls -l|grep -v grep|grep HotfixConfig_$num |wc -l ");
$a5=$ssh2->capture("cd /www/doc/1/android/HotfixConfig/ && ls -l|grep -v grep |grep HotfixConfig_$num |wc -l ");
$a6=$ssh2->capture("cd /www/doc/1/wp8/HotfixConfig/ && ls -l|grep -v grep|grep HotfixConfig_$num |wc -l ");

if($a1==1 && $a2==1 && $a3==1 && $a4==1 && $a5==1 && $a6==1)
{
print "\n\n\n";
print "====================================\n";
print "||版本号为$num的静态数据传输完毕！||\n";
print "====================================\n";
print "\n\n\n";
}
else{
	print "\n\n\n";
	print "================================\n";
	print "||静态数据传输未成功，请检查！||\n";
	print "================================\n";
	print "\n\n\n";
	exit;
	}	
				
			
}



&total();
