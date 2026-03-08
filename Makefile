DEST = systemmgr@10.0.1.150:web/gigaset/info/

FILES = index.php menu.php weather.php .htaccess

publish:
	scp $(FILES) $(DEST)
