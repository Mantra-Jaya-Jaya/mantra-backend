#!/bin/bash
# ==========================================
# MIGRATION SCRIPT: Normalisasi StatusPesanan
# ==========================================

set -e  # Exit on error

echo "🚀 Starting StatusPesanan Normalization Migration"
echo "=================================================="
echo ""

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
else
    echo "❌ Error: .env file not found!"
    exit 1
fi

# Check required environment variables
if [ -z "$DB_HOST" ] || [ -z "$DB_PORT" ] || [ -z "$DB_NAME" ] || [ -z "$DB_USER" ] || [ -z "$DB_PASSWORD" ]; then
    echo "❌ Error: Missing database environment variables!"
    exit 1
fi

DB_URL="postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}"

echo "📍 Database: ${DB_HOST}:${DB_PORT}/${DB_NAME}"
echo ""

# Step 1: Backup database
echo "📦 Step 1: Creating database backup..."
BACKUP_FILE="backup_before_status_normalization_$(date +%Y%m%d_%H%M%S).sql"
PGPASSWORD=${DB_PASSWORD} pg_dump -h ${DB_HOST} -p ${DB_PORT} -U ${DB_USER} -d ${DB_NAME} > ${BACKUP_FILE}
if [ $? -eq 0 ]; then
    echo "✅ Backup created: ${BACKUP_FILE}"
else
    echo "❌ Backup failed!"
    exit 1
fi
echo ""

# Step 2: Run data migration SQL
echo "🗃️  Step 2: Running data migration..."
PGPASSWORD=${DB_PASSWORD} psql -h ${DB_HOST} -p ${DB_PORT} -U ${DB_USER} -d ${DB_NAME} -f migrations/001_status_pesanan_data.sql
if [ $? -eq 0 ]; then
    echo "✅ Data migration completed"
else
    echo "❌ Data migration failed! Check logs above."
    echo "💡 Rollback: Restore from ${BACKUP_FILE}"
    exit 1
fi
echo ""

# Step 3: Verify data migration
echo "🔍 Step 3: Verifying data migration..."
UNMIGRATED=$(PGPASSWORD=${DB_PASSWORD} psql -h ${DB_HOST} -p ${DB_PORT} -U ${DB_USER} -d ${DB_NAME} -t -c "
SELECT COUNT(*) FROM pesanan WHERE status_pesanan_id IS NULL AND status_pesanan IS NOT NULL;
" | tr -d ' ')

if [ "$UNMIGRATED" != "0" ]; then
    echo "⚠️  Warning: ${UNMIGRATED} records failed to migrate!"
    echo "💡 Check the data manually before proceeding"
    read -p "Continue with Atlas schema migration? (y/n): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
else
    echo "✅ All records migrated successfully"
fi
echo ""

# Step 4: Preview Atlas schema changes
echo "👀 Step 4: Previewing Atlas schema changes (dry-run)..."
make db-plan
echo ""

read -p "Apply Atlas schema changes? (y/n): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "❌ Migration cancelled by user"
    exit 1
fi
echo ""

# Step 5: Apply Atlas schema changes
echo "⚙️  Step 5: Applying Atlas schema changes..."
make db-apply
if [ $? -eq 0 ]; then
    echo "✅ Schema migration completed"
else
    echo "❌ Schema migration failed!"
    echo "💡 Rollback: Restore from ${BACKUP_FILE}"
    exit 1
fi
echo ""

# Step 6: Final verification
echo "🎯 Step 6: Final verification..."
echo ""
echo "Checking status_pesanan table:"
PGPASSWORD=${DB_PASSWORD} psql -h ${DB_HOST} -p ${DB_PORT} -U ${DB_USER} -d ${DB_NAME} -c "
SELECT * FROM status_pesanan ORDER BY id_status_pesanan;
"
echo ""

echo "Checking pesanan table sample:"
PGPASSWORD=${DB_PASSWORD} psql -h ${DB_HOST} -p ${DB_PORT} -U ${DB_USER} -d ${DB_NAME} -c "
SELECT id_pesanan, status_pesanan_id, nama_status
FROM pesanan
JOIN status_pesanan ON status_pesanan_id = id_status_pesanan
ORDER BY id_pesanan DESC
LIMIT 10;
"
echo ""

echo "=================================================="
echo "🎉 Migration completed successfully!"
echo ""
echo "📋 Summary:"
echo "  - Backup created: ${BACKUP_FILE}"
echo "  - Data migrated: All existing records"
echo "  - Schema updated: Via Atlas"
echo ""
echo "⚠️  Next steps:"
echo "  1. Test all API endpoints"
echo "  2. Verify frontend functionality"
echo "  3. Update API documentation"
echo "  4. Monitor error logs"
echo ""
echo "🔄 Rollback command (if needed):"
echo "   PGPASSWORD=${DB_PASSWORD} psql -h ${DB_HOST} -p ${DB_PORT} -U ${DB_USER} -d ${DB_NAME} < ${BACKUP_FILE}"
echo ""
