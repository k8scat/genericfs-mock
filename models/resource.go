package models

import (
	"database/sql"
)

type Resource struct {
	UUID          string `db:"uuid" json:"uuid"`
	ReferenceType int    `db:"reference_type" json:"reference_type"`
	ReferenceID   string `db:"reference_id" json:"reference_id"`
	TeamUUID      string `db:"team_uuid" json:"team_uuid"`
	ProjectUUID   string `db:"project_uuid" json:"project_uuid"`
	OwnerUUID     string `db:"owner_uuid" json:"owner_uuid"`
	Modifier      string `db:"modifier" json:"modifier"`
	Type          int    `db:"type" json:"type"`
	Source        int    `db:"source" json:"source"`
	ExtID         string `db:"ext_id" json:"ext_id"`
	Name          string `db:"name" json:"name"`
	Status        int    `db:"status" json:"status"`
	CreateTime    int64  `db:"create_time" json:"create_time"` // time.millisecond
	Description   string `db:"description" json:"description"`
	ModifyTime    int64  `db:"modify_time" json:"modify_time"` // time.milliseconds

	CallbackURL  string `db:"callback_url" json:"callback_url"`
	CallbackBody string `db:"callback_body" json:"callback_body"`

	IsPublic bool `db:"is_public" json:"-"`
}

func GetResourceByUUID(uuid string) (*Resource, error) {
	query := "SELECT uuid, reference_type, reference_id, team_uuid, project_uuid, owner_uuid, modifier, "
	query += "type, source, ext_id, name, create_time, modify_time, status, description, callback_url, callback_body, is_public "
	query += "FROM resource WHERE uuid=?;"
	row := DB.QueryRow(query, uuid)
	resource := new(Resource)
	err := row.Scan(&resource.UUID,
		&resource.ReferenceType,
		&resource.ReferenceID,
		&resource.TeamUUID,
		&resource.ProjectUUID,
		&resource.OwnerUUID,
		&resource.Modifier,
		&resource.Type,
		&resource.Source,
		&resource.ExtID,
		&resource.Name,
		&resource.CreateTime,
		&resource.ModifyTime,
		&resource.Status,
		&resource.Description,
		&resource.CallbackURL,
		&resource.CallbackBody,
		&resource.IsPublic)
	return resource, err
}

func GetResourceByHash(hash string) (*Resource, error) {
	query := "SELECT uuid, reference_type, reference_id, team_uuid, project_uuid, owner_uuid, modifier, "
	query += "type, source, ext_id, name, create_time, modify_time, status, description, callback_url, callback_body, is_public "
	query += "FROM resource WHERE ext_id=?;"
	row := DB.QueryRow(query, hash)
	resource := new(Resource)
	err := row.Scan(&resource.UUID,
		&resource.ReferenceType,
		&resource.ReferenceID,
		&resource.TeamUUID,
		&resource.ProjectUUID,
		&resource.OwnerUUID,
		&resource.Modifier,
		&resource.Type,
		&resource.Source,
		&resource.ExtID,
		&resource.Name,
		&resource.CreateTime,
		&resource.ModifyTime,
		&resource.Status,
		&resource.Description,
		&resource.CallbackURL,
		&resource.CallbackBody,
		&resource.IsPublic)
	return resource, err
}

func AddResource(resource *Resource) error {
	query := "INSERT INTO resource(uuid, reference_type, reference_id, team_uuid, project_uuid, owner_uuid, modifier, "
	query += "type, source, ext_id, name, create_time, modify_time, status, description, callback_url, callback_body, is_public) "
	query += "VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);"
	return Transact(func(tx *sql.Tx) error {
		_, err := tx.Exec(query,
			resource.UUID,
			resource.ReferenceType,
			resource.ReferenceID,
			resource.TeamUUID,
			resource.ProjectUUID,
			resource.OwnerUUID,
			resource.Modifier,
			resource.Type,
			resource.Source,
			resource.ExtID,
			resource.Name,
			resource.CreateTime,
			resource.ModifyTime,
			resource.Status,
			resource.Description,
			resource.CallbackURL,
			resource.CallbackBody,
			resource.IsPublic)
		return err
	})
}

func UpdateResource(resource *Resource) error {
	query := "UPDATE resource SET reference_type=?, reference_id=?, team_uuid=?, project_uuid=?, modifier=?, "
	query += "type=?, ext_id=?, name=?, create_time=?, modify_time=?, status=?, description=? WHERE uuid=?;"
	return Transact(func(tx *sql.Tx) error {
		_, err := tx.Exec(query,
			resource.ReferenceType,
			resource.ReferenceID,
			resource.TeamUUID,
			resource.ProjectUUID,
			resource.Modifier,
			resource.Type,
			resource.ExtID,
			resource.Name,
			resource.CreateTime,
			resource.ModifyTime,
			resource.Status,
			resource.Description,
			resource.UUID)
		return err
	})
}

func IsFileExisted(hash string) (bool, error) {
	query := "SELECT COUNT(1) FROM resource WHERE ext_id=?;"
	var count int
	row := DB.QueryRow(query, hash)
	err := row.Scan(&count)
	return count > 0, err
}
