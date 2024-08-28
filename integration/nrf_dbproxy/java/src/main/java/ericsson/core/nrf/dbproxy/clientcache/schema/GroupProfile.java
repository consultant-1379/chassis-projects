package ericsson.core.nrf.dbproxy.clientcache.schema;

import com.google.protobuf.ByteString;
import ericsson.core.nrf.dbproxy.common.Code;
import java.util.HashMap;

public class GroupProfile {

  private String group_profile_id;
  private HashMap nf_type;
  private HashMap group_id;
  private int profile_type;
  private Long supi_version;
  private ByteString raw_group_profile;

  public GroupProfile() {
    this.group_profile_id = "";
    this.nf_type = new HashMap();
    this.group_id = new HashMap();
    this.profile_type = Code.PROFILE_TYPE_EMPTY;
    this.supi_version = 0L;
    this.raw_group_profile = ByteString.EMPTY;
  }

  public String getGroupProfileID() {
    return group_profile_id;
  }

  public void setGroupProfileID(String groupProfileId) {
    this.group_profile_id = groupProfileId;
  }


  public void addNFType(String nfType) {
    this.nf_type.put(nfType, 1);
  }

  public void addGroupID(String groupId) {
    this.group_id.put(groupId, 1);
  }

  public void setProfileType(int profileType) {
    this.profile_type = profileType;
  }

  public Long getSupiVersion() {
    return this.supi_version;
  }

  public void setSupiVersion(Long supiVersion) {
    this.supi_version = supiVersion;
  }

  public ByteString getData() {
    return raw_group_profile;
  }

  public void setData(ByteString rawGroupProfile) {
    this.raw_group_profile = rawGroupProfile;
  }

  public String toString() {
    String nfTypeStr = "";
    if (!nf_type.isEmpty()) {
      nfTypeStr = nf_type.toString();
    }
    String groupIdStr = "";
    if (!group_id.isEmpty()) {
      groupIdStr = group_id.toString();
    }
    return "GroupProfile{" +
        "group_profile_id=" + group_profile_id +
        ", nf_type=" + nfTypeStr +
        ", group_id=" + groupIdStr +
        ", profile_type=" + profile_type +
        ", supi_version=" + supi_version.toString() +
        ", raw_group_profile=" + raw_group_profile.toString() +
        '}';
  }
}
