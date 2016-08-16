mkdir d:\db
mkdir d:\db\data
mkdir d:\db\log
del d:\db\data\mongod.lock
..\..\tools\mongodb\mongod.exe --dbpath d:\db\data --logpath d:\db\mongodb.log --logappend
