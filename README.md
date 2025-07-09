# 🌀 mysync – A Mini Rsync-like File Synchronizer in Go

**mysync** is a lightweight file synchronization tool built in Go that allows syncing directories either locally or to a remote machine over a custom TCP protocol — similar to a basic version of `rsync`.

---

## 🚀 Features

- ✅ Syncs directories and files from source to destination
- ✅ Supports both local and remote sync
- ✅ Compares files using size and modification time
- ✅ Automatically creates missing directories at destination
- ✅ Efficient skipping of unchanged files
- ✅ `--mirror` mode: deletes files/dirs in destination not present in source
- ✅ `--dry-run` mode: previews changes without copying
- ✅ Persistent TCP connection to reduce overhead
- ✅ Simple custom protocol over TCP (no SSH or SCP needed)

---

## ⚙️ Usage

### 🔧 Build

```bash
go build -o mysync
```

## 🔄 Local Sync

```bash
./mysync ./src ./dst
```

## 🌐 Remote Sync
Start the server on the remote machine:

```bash
go run server.go
```

On the client:

```bash
./mysync --remote 192.168.0.101:8080 ./src /remote/absolute/path
```

## 🧹 Mirror Mode

Deletes extra files in the destination not present in the source. (Currently supported in local sync)

```bash
./mysync --mirror ./src ./dst
```

## 👀 Dry Run
Shows actions without actually copying/deleting files. (Currently supported in local sync)

```bash
./mysync --dry-run ./src ./dst
```