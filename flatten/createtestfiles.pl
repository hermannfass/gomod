#!/usr/bin/env perl

use 5.010;
use utf8;
use File::Spec;

say "Running Perl5 version $^V";

my $top_path = './testfiles';
my $levels_to_do = 3;
my @dir_names = ('A'..'D');

unless (-d $top_path) {
	say "Directory $top_path gets created.";
	mkdir($top_path, 0755) or die "Cannot create $top_path. $!";
}

&create_dirs($top_path, \@dir_names);

sub create_dirs {
	my $parent_path = shift;
	my $dns = shift;
	state $levels_done = 0;
	foreach my $dn (@$dns) {
		my $dp = File::Spec->catfile($parent_path, $dn);
		if (-e $dp) {
			say "Skipping existing directory: $dp";
			next;
		}
		say "Creating directory: $dp";
		mkdir $dp or die "Could not create directory $dp. $!";
		# To do / implement:
		&populate_dir($dp);
		$levels_done++;
		if ($levels_done < $levels_to_do) {
			my $next_parent_path = File::Spec->catfile($parent_path, $dn);
			say "Creating in $next_parent_path:";
			map { say " - $_" } @$dns;
			&create_dirs($next_parent_path, $dns);		
		}
		$levels_done--;
		say "$levels_done levels done.";
	}
}

sub populate_dir {
	my $dir_path = shift;
	say "Creating files in: $dir_path";
	foreach my $fn ('a'..'c') {
		my $fp = File::Spec->catfile($dir_path, "$fn.txt");
		say "  Creating file: $fp";
		my $content = "File: $fn\nLocation: $dir_path\nFQFN: $fp\n";
		open my $fh, '>', $fp or die "Cannot open $fp. $!";
		print $fh $content;
		close $fh;
	}
}

