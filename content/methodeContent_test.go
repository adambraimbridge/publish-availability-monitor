package content

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

//testExternalValidationEndpoints
const testExtValEndpoint = "http://transformer/content-transorm/"

func TestIsEomfileValid_InvalidContentType(t *testing.T) {
	if eomfileWithInvalidContentType.IsValid(testExtValEndpoint, "", "") {
		t.Error("Eomfile with invalid content marked as valid")
	}
}

func TestIsEomfileValid_InvalidUUID(t *testing.T) {
	if eomfileWithInvalidUUID.IsValid(testExtValEndpoint, "", "") {
		t.Error("Eomfile with invalid UUID marked as valid")
	}
}

func TestIsEomfileValid_InvalidSourceCode(t *testing.T) {
	if unsupportedSourceCodeCompoundStory.IsValid(testExtValEndpoint, "", "") {
		t.Error("Eomfile with unsupported source code marked as valid")
	}
}

func TestIsEomfileValid_ValidImage(t *testing.T) {
	if !validImage.IsValid(testExtValEndpoint, "", "") {
		t.Error("Valid Image marked as invalid!")
	}
}

func TestIsEomfileValid_ValidList(t *testing.T) {
	if !validList.IsValid(testExtValEndpoint, "", "") {
		t.Error("Valid List marked as invalid!")
	}
}

func TestIsEomfileValid_ExternalValidationTrue_ValidCompoundStory(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//return OK
	}))
	defer ts.Close()
	if !validCompoundStory.IsValid(ts.URL, "", "") {
		t.Error("Valid CompoundStory marked as invalid!")
	}
}

func TestIsEomfileValid_ExternalValidationFalse_InvalidCompountStory(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(ImATeapot)
	}))
	defer ts.Close()
	if validCompoundStory.IsValid(ts.URL, "", "") {
		t.Error("Valid CompoundStory regarded as invalid by external validation marked as valid!")
	}
}

func TestIsEomfileValid_ExternalValidationTrue_ValidStory(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//return OK
	}))

	if !validStory.IsValid(ts.URL, "", "") {
		t.Error("Valid Story marked as invalid!")
	}
}

func TestIsEomfileValid_ExternalValidationFalse_ValidStory(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(ImATeapot)
	}))

	if validStory.IsValid(ts.URL, "", "") {
		t.Error("Valid Story regarded as invalid by external validation marked as valid!")
	}
}

func TestHasTitle_ValidTitle(t *testing.T) {
	if !hasTitle(eomfileWithTitle) {
		t.Error("Eom File with title marked as invalid!")
	}
}

func TestHasTitle_invalidTitle(t *testing.T) {
	if hasTitle(eomfileWithoutTitle) {
		t.Error("Eom File without title marked as valid!")
	}
}

func TestIsSupportedChannel_CompoundStoryWebChannel(t *testing.T) {
	var validChannelEomFile = EomFile{
		UUID:             validUUID,
		Type:             "EOM::CompoundStory",
		Value:            "bar",
		Attributes:       "attributes",
		SystemAttributes: systemAttributesWebChannel,
	}

	if !isSupportedChannel(validChannelEomFile) {
		t.Error("Eom File with valid channel  marked as invalid!")
	}
}

func TestIsSupportedChannel_CompoundStoryFTChannel(t *testing.T) {
	var invalidChannelEomFile = EomFile{
		UUID:             validUUID,
		Type:             "EOM::CompoundStory",
		Value:            "bar",
		Attributes:       "attributes",
		SystemAttributes: systemAttributesFTChannel,
	}

	if isSupportedChannel(invalidChannelEomFile) {
		t.Error("Eom File with invalid channel  marked as valid!")
	}
}

func TestIsSupportedChannel_CompoundStoryInvalidChannel(t *testing.T) {
	var invalidChannelEomFile = EomFile{
		UUID:             validUUID,
		Type:             "EOM::CompoundStory",
		Value:            "bar",
		Attributes:       "attributes",
		SystemAttributes: systemAttributesInvalidChannel,
	}

	if isSupportedChannel(invalidChannelEomFile) {
		t.Error("Eom File with invalid channel  marked as valid!")
	}
}

func TestIsSupportedChannel_StoryWebChannel(t *testing.T) {
	var validChannelEomFile = EomFile{
		UUID:             validUUID,
		Type:             "EOM::Story",
		Value:            "bar",
		Attributes:       "attributes",
		SystemAttributes: systemAttributesWebChannel,
	}

	if !isSupportedChannel(validChannelEomFile) {
		t.Error("Eom File with valid channel  marked as invalid!")
	}
}

func TestIsSupportedChannel_StoryFTChannel(t *testing.T) {
	var validChannelEomFile = EomFile{
		UUID:             validUUID,
		Type:             "EOM::Story",
		Value:            "bar",
		Attributes:       "attributes",
		SystemAttributes: systemAttributesFTChannel,
	}

	if !isSupportedChannel(validChannelEomFile) {
		t.Error("Eom File with valid channel  marked as invalid!")
	}
}

func TestIsSupportedChannel_StoryInvalidChannel(t *testing.T) {
	var invalidChannelEomFile = EomFile{
		UUID:             validUUID,
		Type:             "EOM::Story",
		Value:            "bar",
		Attributes:       "attributes",
		SystemAttributes: systemAttributesInvalidChannel,
	}

	if isSupportedChannel(invalidChannelEomFile) {
		t.Error("Eom File with invalid channel  marked as valid!")
	}
}

func TestIsSupportedFileType_SupportedType(t *testing.T) {
	if !isSupportedFileType(supportedEomFile) {
		t.Error("Eom File with supported filetype marked as invalid!")
	}
}

func TestIsSupportedFileType_UnupportedType(t *testing.T) {
	if isSupportedFileType(unsupportedEomFile) {
		t.Error("Eom File with unsupported filetype marked as valid!")
	}
}

func TestIsSupportedSourceCode_CompoundStory_SupportedCode(t *testing.T) {
	if !isSupportedCompoundStorySourceCode(supportedSourceCodeCompoundStory) {
		t.Error("Compound story with supported source code marked as invalid!")
	}
}

func TestIsSupportedSourceCode_ContentplaceHolderCompoundStory_SupportedCode(t *testing.T) {
	if !isSupportedCompoundStorySourceCode(contentplaceHolderCompoundStory) {
		t.Error("Compound story with supported source code marked as invalid!")
	}
}

func TestIsSupportedSourceCodes_Story_SupportedCode(t *testing.T) {
	if !isSupportedStorySourceCode(supportedSourceCodeStory) {
		t.Error("Story with supported source code marked as invalid!")
	}
}

func TestIsSupportedSourceCode_CompoundStory_UnsupportedCode(t *testing.T) {
	if isSupportedCompoundStorySourceCode(unsupportedSourceCodeCompoundStory) {
		t.Error("Story with unsupported source code marked as valid!")
	}
}

func TestIsSupportedSourceCode_Story_UnsupportedCode(t *testing.T) {
	if isSupportedStorySourceCode(unsupportedSourceCodeStory) {
		t.Error("Compound story with unsupported source code marked as valid!")
	}
}

func TestIsImageValid_ImageValid(t *testing.T) {
	if !isImageValid(validImageEomFile) {
		t.Error("Valid Image EOMFile marked as invalid!")
	}
}

func TestIsImageValid_ImageInvalid(t *testing.T) {
	if isImageValid(invalidImageEomFile) {
		t.Error("Invalid Image EOMFile marked as valid!")
	}
}

func TestIsListValid_ListValid(t *testing.T) {
	if !isListValid(validListEomFile) {
		t.Error("Valid List EOMFile marked as invalid!")
	}
}

func TestIsListValid_ListInvalid(t *testing.T) {
	if isListValid(invalidListEomFile) {
		t.Error("Invalid List EOMFile marked as valid!")
	}
}

var eomfileWithInvalidContentType = EomFile{
	UUID:             validUUID,
	Type:             "FOOBAR",
	Value:            "/9j/4QAYRXhpZgAASUkqAAgAAAAAAAAAAAAAAP/sABFEdWNr",
	Attributes:       "attributes",
	SystemAttributes: "systemAttributes",
}
var eomfileWithInvalidUUID = EomFile{
	UUID:             invalidUUID,
	Type:             "Image",
	Value:            "/9j/4QAYRXhpZgAASUkqAAgAAAAAAAAAAAAAAP/sABFEdWNr",
	Attributes:       "attributes",
	SystemAttributes: "systemAttributes",
}

var validImage = EomFile{
	UUID:             validUUID,
	Type:             "Image",
	Value:            "/9j/4QAYRXhpZgAASUkqAAgAAAAAAAAAAAAAAP/sABFEdWNr",
	Attributes:       "attributes",
	SystemAttributes: "systemAttributes",
}

var validList = EomFile{
	UUID:             validUUID,
	Type:             "EOM::WebContainer",
	Value:            "bar",
	Attributes:       validListAttributes,
	SystemAttributes: "system attributes",
}

var validCompoundStory = EomFile{
	UUID:             validUUID,
	Type:             "EOM::CompoundStory",
	Value:            contentWithHeadline,
	Attributes:       validFileTypeAttributes,
	SystemAttributes: systemAttributesWebChannel,
}

var validStory = EomFile{
	UUID:             validUUID,
	Type:             "EOM::Story",
	Value:            contentWithHeadline,
	Attributes:       supportedSourceCodeAttributes,
	SystemAttributes: systemAttributesWebChannel,
}

var eomfileWithTitle = EomFile{
	UUID:             validUUID,
	Type:             "EOM::CompoundStory",
	Value:            contentWithHeadline,
	Attributes:       "attributes",
	SystemAttributes: "systemAttributes",
}

var eomfileWithoutTitle = EomFile{
	UUID:             validUUID,
	Type:             "EOM::CompoundStory",
	Value:            contentWithoutHeadline,
	Attributes:       "attributes",
	SystemAttributes: "systemAttributes",
}

var supportedSourceCodeCompoundStory = EomFile{
	UUID:             validUUID,
	Type:             "EOM::CompoundStory",
	Value:            "bar",
	Attributes:       supportedSourceCodeAttributes,
	SystemAttributes: "systemAttributes",
}

var contentplaceHolderCompoundStory = EomFile{
	UUID:             validUUID,
	Type:             "EOM::CompoundStory",
	Value:            "bar",
	Attributes:       supportedSourceCodeAttributesContentPlaceholder,
	SystemAttributes: "systemAttributes",
}

var supportedSourceCodeStory = EomFile{
	UUID:             validUUID,
	Type:             "EOM::Story",
	Value:            "value",
	Attributes:       supportedSourceCodeAttributes,
	SystemAttributes: "systemAttributes",
}

var unsupportedSourceCodeCompoundStory = EomFile{
	UUID:             validUUID,
	Type:             "EOM::CompoundStory",
	Value:            "bar",
	Attributes:       unsupportedSourceCodeAttributes,
	SystemAttributes: "systemAttributes",
}

var unsupportedSourceCodeStory = EomFile{
	UUID:             validUUID,
	Type:             "EOM::Story",
	Value:            "bar",
	Attributes:       unsupportedSourceCodeAttributes,
	SystemAttributes: "systemAttributes",
}
var supportedEomFile = EomFile{
	UUID:             validUUID,
	Type:             "EOM::CompoundStory",
	Value:            "bar",
	Attributes:       validFileTypeAttributes,
	SystemAttributes: "system attributes",
}

var unsupportedEomFile = EomFile{
	UUID:             validUUID,
	Type:             "EOM::CompoundStory",
	Value:            "bar",
	Attributes:       invalidFileTypeAttributes,
	SystemAttributes: "system attributes",
}

var validListEomFile = EomFile{
	UUID:             validUUID,
	Type:             "EOM::WebContainer",
	Value:            "bar",
	Attributes:       validListAttributes,
	SystemAttributes: "system attributes",
}

var invalidListEomFile = EomFile{
	UUID:             validUUID,
	Type:             "EOM::WebContainer",
	Value:            "bar",
	Attributes:       invalidListAttributes,
	SystemAttributes: "system attributes",
}

var validImageEomFile = EomFile{
	UUID:             validUUID,
	Type:             "Image",
	Value:            "/9j/4QAYRXhpZgAASUkqAAgAAAAAAAAAAAAAAP/sABFEdWNr",
	Attributes:       "attributes",
	SystemAttributes: "system attributes",
}

var invalidImageEomFile = EomFile{
	UUID:             validUUID,
	Type:             "Image",
	Value:            "",
	Attributes:       "attributes",
	SystemAttributes: "system attributes",
}

func TestIsMarkedDeleted_CompoundStory_True(t *testing.T) {
	if !compoundStoryMarkedDeletedTrue.IsMarkedDeleted() {
		t.Error("Expected True, the compound story IS marked deleted")
	}
}

func TestIsMarkedDeleted_CompoundStory_False(t *testing.T) {
	if compoundStoryMarkedDeletedFalse.IsMarkedDeleted() {
		t.Error("Expected False, the compound story IS NOT marked deleted")
	}
}

func TestIsMarkedDeleted_Story_True(t *testing.T) {
	if !storyMarkedDeletedTrue.IsMarkedDeleted() {
		t.Error("Expected True, the story IS marked deleted")
	}
}

func TestIsMarkedDeleted_Story_False(t *testing.T) {
	if storyMarkedDeletedFalse.IsMarkedDeleted() {
		t.Error("Expected False, the story IS NOT marked deleted")
	}
}

func TestIsMarkedDeleted_Image(t *testing.T) {
	if imageEomFile.IsMarkedDeleted() {
		t.Error("Expected False, the image IS NOT marked deleted")
	}
}

func TestIsMarkedDeleted_WebContainer(t *testing.T) {
	if webContainerEomFile.IsMarkedDeleted() {
		t.Error("Expected False, the webContainer IS NOT marked deleted")
	}
}

var compoundStoryMarkedDeletedTrue = EomFile{
	UUID:             UUIDString,
	Type:             "EOM::CompoundStory",
	Value:            content,
	Attributes:       attributesMarkedDeletedTrue,
	SystemAttributes: systemAttributes,
}

var storyMarkedDeletedTrue = EomFile{
	UUID:             UUIDString,
	Type:             "EOM::Story",
	Value:            content,
	Attributes:       attributesMarkedDeletedTrue,
	SystemAttributes: systemAttributes,
}

var imageEomFile = EomFile{
	UUID:             UUIDString,
	Type:             "Image",
	Value:            "image bytes",
	Attributes:       "fooAttributes",
	SystemAttributes: "barsystemAttributes",
}

var compoundStoryMarkedDeletedFalse = EomFile{
	UUID:             UUIDString,
	Type:             "EOM::CompoundStory",
	Value:            content,
	Attributes:       attributesMarkedDeletedFalse,
	SystemAttributes: systemAttributes,
}

var storyMarkedDeletedFalse = EomFile{
	UUID:             UUIDString,
	Type:             "EOM::Story",
	Value:            content,
	Attributes:       attributesMarkedDeletedFalse,
	SystemAttributes: systemAttributes,
}

var webContainerEomFile = EomFile{
	UUID:             UUIDString,
	Type:             "EOM::WebContainer",
	Value:            "list bytes",
	Attributes:       "fooAttributes",
	SystemAttributes: "barSystemAttributes",
}

const UUIDString = "e28b12f7-9796-3331-b030-05082f0b8157"
const systemAttributes = "<props><productInfo><name>FTcom</name><issueDate>20150915</issueDate></productInfo><workFolder>/FT/WorldNews</workFolder><subFolder>UKNews</subFolder></props>"
const content = "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0idXRmLTgiPz4NCjwhRE9DVFlQRSBkb2MgU1lTVEVNICIvU3lzQ29uZmlnL1J1bGVzL2Z0cHNpLmR0ZCI+DQo8P0VNLWR0ZEV4dCAvU3lzQ29uZmlnL1J1bGVzL2Z0cHNpL2Z0cHNpLmR0eD8+DQo8P0VNLXRlbXBsYXRlTmFtZSAvU3lzQ29uZmlnL1RlbXBsYXRlcy9GVC9CYXNlLVN0b3J5LnhtbD8+DQo8P3htbC1mb3JtVGVtcGxhdGUgL1N5c0NvbmZpZy9UZW1wbGF0ZXMvRlQvQmFzZS1TdG9yeS54cHQ/Pg0KPD94bWwtc3R5bGVzaGVldCB0eXBlPSJ0ZXh0L2NzcyIgaHJlZj0iL1N5c0NvbmZpZy9SdWxlcy9mdHBzaS9GVC9tYWlucmVwLmNzcyI/Pg0KPGRvYyB4bWw6bGFuZz0iZW4tdWsiPjxsZWFkIGlkPSJVMzIwMDExNzk2MDc0NzhGQkUiPjxsZWFkLWhlYWRsaW5lIGlkPSJVMzIwMDExNzk2MDc0NzhhMkQiPjxuaWQtdGl0bGU+PGxuPjw/RU0tZHVtbXlUZXh0IFtJbnNlcnQgbmV3cyBpbiBkZXB0aCB0aXRsZSBoZXJlXT8+DQo8L2xuPg0KPC9uaWQtdGl0bGU+DQo8aW4tZGVwdGgtbmF2LXRpdGxlPjxsbj48P0VNLWR1bW15VGV4dCBbSW5zZXJ0IGluIGRlcHRoIG5hdiB0aXRsZSBoZXJlXT8+DQo8L2xuPg0KPC9pbi1kZXB0aC1uYXYtdGl0bGU+DQo8aGVhZGxpbmUgaWQ9IlUzMjAwMTE3OTYwNzQ3OHNtQyI+PGxuPlRoaXMgaXMgYSBoZWFkbGluZQ0KPC9sbj4NCjwvaGVhZGxpbmU+DQo8c2t5Ym94LWhlYWRsaW5lIGlkPSJVMzIwMDExNzk2MDc0Nzg1a0YiPjxsbj48P0VNLWR1bW15VGV4dCBbU2t5Ym94IGhlYWRsaW5lIGhlcmVdPz4NCjwvbG4+DQo8L3NreWJveC1oZWFkbGluZT4NCjx0cmlwbGV0LWhlYWRsaW5lPjxsbj48P0VNLWR1bW15VGV4dCBbVHJpcGxldCBoZWFkbGluZSBoZXJlXT8+DQo8L2xuPg0KPC90cmlwbGV0LWhlYWRsaW5lPg0KPHByb21vYm94LXRpdGxlPjxsbj48P0VNLWR1bW15VGV4dCBbUHJvbW9ib3ggdGl0bGUgaGVyZV0/Pg0KPC9sbj4NCjwvcHJvbW9ib3gtdGl0bGU+DQo8cHJvbW9ib3gtaGVhZGxpbmU+PGxuPjw/RU0tZHVtbXlUZXh0IFtQcm9tb2JveCBoZWFkbGluZSBoZXJlXT8+DQo8L2xuPg0KPC9wcm9tb2JveC1oZWFkbGluZT4NCjxlZGl0b3ItY2hvaWNlLWhlYWRsaW5lIGlkPSJVMzIwMDExNzk2MDc0NzhGVkUiPjxsbj48P0VNLWR1bW15VGV4dCBbU3RvcnkgcGFja2FnZSBoZWFkbGluZV0/Pg0KPC9sbj4NCjwvZWRpdG9yLWNob2ljZS1oZWFkbGluZT4NCjxuYXYtY29sbGVjdGlvbi1oZWFkbGluZT48bG4+PD9FTS1kdW1teVRleHQgW05hdiBjb2xsZWN0aW9uIGhlYWRsaW5lXT8+DQo8L2xuPg0KPC9uYXYtY29sbGVjdGlvbi1oZWFkbGluZT4NCjxpbi1kZXB0aC1uYXYtaGVhZGxpbmU+PGxuPjw/RU0tZHVtbXlUZXh0IFtJbiBkZXB0aCBuYXYgaGVhZGxpbmVdPz4NCjwvbG4+DQo8L2luLWRlcHRoLW5hdi1oZWFkbGluZT4NCjwvbGVhZC1oZWFkbGluZT4NCjx3ZWItaW5kZXgtaGVhZGxpbmUgaWQ9IlUzMjAwMTE3OTYwNzQ3OHRQRyI+PGxuPjw/RU0tZHVtbXlUZXh0IFtTaG9ydCBoZWFkbGluZV0/Pg0KPC9sbj4NCjwvd2ViLWluZGV4LWhlYWRsaW5lPg0KPHdlYi1zdGFuZC1maXJzdCBpZD0iVTMyMDAxMTc5NjA3NDc4bmVFIj48cD48P0VNLWR1bW15VGV4dCBbTG9uZyBzdGFuZGZpcnN0XT8+DQo8L3A+DQo8L3dlYi1zdGFuZC1maXJzdD4NCjx3ZWItc3ViaGVhZCBpZD0iVTMyMDAxMTc5NjA3NDc4dGpHIj48cD48P0VNLWR1bW15VGV4dCBbU2hvcnQgc3RhbmRmaXJzdF0/Pg0KPC9wPg0KPC93ZWItc3ViaGVhZD4NCjxlZGl0b3ItY2hvaWNlIGlkPSJVMzIwMDExNzk2MDc0NzhRWUUiPjw/RU0tZHVtbXlUZXh0IFtTdG9yeSBwYWNrYWdlIGh5cGVybGlua10/Pg0KPC9lZGl0b3ItY2hvaWNlPg0KPGxlYWQtdGV4dCBpZD0iVTMyMDAxMTc5NjA3NDc4T0VEIj48bGVhZC1ib2R5PjxwPjw/RU0tZHVtbXlUZXh0IFtJbnNlcnQgbGVhZCBib2R5IHRleHQgaGVyZSAtIG1pbiAxMzAgY2hhcnMsIG1heCAxNTAgY2hhcnNdPz4NCjwvcD4NCjwvbGVhZC1ib2R5Pg0KPHRyaXBsZXQtbGVhZC1ib2R5PjxwPjw/RU0tZHVtbXlUZXh0IFtJbnNlcnQgdHJpcGxldCBsZWFkIGJvZHkgYm9keSBoZXJlXT8+DQo8L3A+DQo8L3RyaXBsZXQtbGVhZC1ib2R5Pg0KPGNvbHVtbmlzdC1sZWFkLWJvZHk+PHA+PD9FTS1kdW1teVRleHQgW0luc2VydCBjb2x1bW5pc3QgbGVhZCBib2R5IGJvZHkgaGVyZV0/Pg0KPC9wPg0KPC9jb2x1bW5pc3QtbGVhZC1ib2R5Pg0KPHNob3J0LWJvZHk+PHA+PD9FTS1kdW1teVRleHQgW0luc2VydCBzaG9ydCBib2R5IGJvZHkgaGVyZV0/Pg0KPC9wPg0KPC9zaG9ydC1ib2R5Pg0KPHNreWJveC1ib2R5PjxwPjw/RU0tZHVtbXlUZXh0IFtJbnNlcnQgc2t5Ym94IGJvZHkgaGVyZV0/Pg0KPC9wPg0KPC9za3lib3gtYm9keT4NCjxwcm9tb2JveC1ib2R5PjxwPjw/RU0tZHVtbXlUZXh0IFtJbnNlcnQgcHJvbW9ib3ggYm9keSBoZXJlXT8+DQo8L3A+DQo8L3Byb21vYm94LWJvZHk+DQo8dHJpcGxldC1zaG9ydC1ib2R5PjxwPjw/RU0tZHVtbXlUZXh0IFtJbnNlcnQgdHJpcGxldCBzaG9ydCBib2R5IGhlcmVdPz4NCjwvcD4NCjwvdHJpcGxldC1zaG9ydC1ib2R5Pg0KPGVkaXRvci1jaG9pY2Utc2hvcnQtbGVhZC1ib2R5PjxwPjw/RU0tZHVtbXlUZXh0IFtJbnNlcnQgZWRpdG9ycyBjaG9pY2Ugc2hvcnQgbGVhZCBib2R5IGhlcmVdPz4NCjwvcD4NCjwvZWRpdG9yLWNob2ljZS1zaG9ydC1sZWFkLWJvZHk+DQo8bmF2LWNvbGxlY3Rpb24tc2hvcnQtbGVhZC1ib2R5PjxwPjw/RU0tZHVtbXlUZXh0IFtJbnNlcnQgbmF2IGNvbGxlY3Rpb24gc2hvcnQgbGVhZCBib2R5IGhlcmVdPz4NCjwvcD4NCjwvbmF2LWNvbGxlY3Rpb24tc2hvcnQtbGVhZC1ib2R5Pg0KPC9sZWFkLXRleHQ+DQo8bGVhZC1pbWFnZXMgaWQ9IlUzMjAwMTE3OTYwNzQ3OEFkQiI+PHdlYi1tYXN0ZXIgaWQ9IlUzMjAwMTE3OTYwNzQ3OFpSRiIvPg0KPHdlYi1za3lib3gtcGljdHVyZS8+DQo8d2ViLWFsdC1waWN0dXJlLz4NCjx3ZWItcG9wdXAtcHJldmlldyB3aWR0aD0iMTY3IiBoZWlnaHQ9Ijk2Ii8+DQo8d2ViLXBvcHVwLz4NCjwvbGVhZC1pbWFnZXM+DQo8aW50ZXJhY3RpdmUtY2hhcnQ+PD9FTS1kdW1teVRleHQgW0ludGVyYWN0aXZlLWNoYXJ0IGxpbmtdPz4NCjwvaW50ZXJhY3RpdmUtY2hhcnQ+DQo8L2xlYWQ+DQo8c3Rvcnk+PGhlYWRibG9jayBpZD0iVTMyMDAxMTc5NjA3NDc4d21IIj48aGVhZGxpbmUgaWQ9IlUzMjAwMTE3OTYwNzQ3OFgwSCI+PGxuPjw/RU0tZHVtbXlUZXh0IFtIZWFkbGluZV0/Pg0KPC9sbj4NCjwvaGVhZGxpbmU+DQo8L2hlYWRibG9jaz4NCjx0ZXh0IGlkPSJVMzIwMDExNzk2MDc0Nzh3WEQiPjxieWxpbmU+PGF1dGhvci1uYW1lPktpcmFuIFN0YWNleSwgRW5lcmd5IENvcnJlc3BvbmRlbnQ8L2F1dGhvci1uYW1lPg0KPC9ieWxpbmU+DQo8Ym9keT48cD5DT1NUIENPTVBBUklTSU9OIFdJVEggT1RIRVIgUkVHSU9OUzwvcD4NCjxwPjwvcD4NCjxwPlRvdmUgU3R1aHIgU2pvYmxvbSBzdG9vZCBpbiBmcm9udCBvZiAzMDAgb2lsIGluZHVzdHJ5IGluc2lkZXJzIGluIEFiZXJkZWVuIGxhc3Qgd2Vlaywgc2xpcHBlZCBpbnRvIHRoZSBsb2NhbCBTY290dGlzaCBkaWFsZWN0LCBhbmQgZ2F2ZSB0aGVtIG9uZSBjbGVhciBtZXNzYWdlOiDigJxLZWVwIHlvdXIgaGVpZC7igJ08L3A+DQo8cD5UaGUgTm9yd2VnaWFuIFN0YXRvaWwgZXhlY3V0aXZlIHdhcyBhZGRyZXNzaW5nIHRoZSBiaWVubmlhbCBPZmZzaG9yZSBFdXJvcGUgY29uZmVyZW5jZSwgZG9taW5hdGVkIHRoaXMgeWVhciBieSB0YWxrIG9mIHRoZSBmYWxsaW5nIG9pbCBwcmljZSBhbmQgdGhlIDxhIGhyZWY9Imh0dHA6Ly93d3cuZnQuY29tL2Ntcy9zLzAvNGYwNDQyZTAtYmNkMi0xMWU0LWE5MTctMDAxNDRmZWFiN2RlLmh0bWwjYXh6ejNsbllyTDN3YiIgdGl0bGU9Ik5vcnRoIFNlYSBvaWw6IFRoYXQgc2lua2luZyBmZWVsaW5nIC0gRlQiPmZ1dHVyZSBvZiBOb3J0aCBTZWEgcHJvZHVjdGlvbjwvYT4uIDwvcD4NCjxwPk92ZXIgdGhlIHBhc3QgMTQgbW9udGhzLCB0aGUgb2lsIHByaWNlIGhhcyBkcm9wcGVkIGZyb20gYWJvdXQgJDExNSBwZXIgYmFycmVsIHRvIGFyb3VuZCAkNTAuIE5vd2hlcmUgaGFzIHRoaXMgYmVlbiBmZWx0IG1vcmUga2Vlbmx5IHRoYW4gaW4gdGhlIE5vcnRoIFNlYSwgRk9SIERFQ0FERVMgYSBwb3dlcmhvdXNlIG9mIHRoZSBCcml0aXNoIGVjb25vbXksIGJ1dCBOT1cgYW4gYXJlYSB3aGVyZSBvaWwgZmllbGRzIGFyZSBnZXR0aW5nIG9sZGVyLCBuZXcgZGlzY292ZXJpZXMgcmFyZXIgYW5kIGNvc3RzIHN0ZWVwZXIuIE5vIHdvbmRlciBNcyBTdHVociBTam9ibG9tIGlzIHRyeWluZyB0byBzdG9wIGhlciBjb2xsZWFndWVzIGZyb20gcGFuaWNraW5nLjwvcD4NCjxwPkJ5IG5lYXJseSBhbnkgbWV0cmljLCBkb2luZyBidXNpbmVzcyBpbiB0aGUgTm9ydGggU2VhIGlzIG1vcmUgZGlmZmljdWx0IHRoYW4gYXQgYW55IG90aGVyIHBlcmlvZCBpbiBpdHMgNTAteWVhciBoaXN0b3J5LiBXaGlsZSB0aGUgb2lsIHByaWNlIGhhcyBiZWVuIGxvd2VyIOKAlCBpdCB0b3VjaGVkICQxMCBhIGJhcnJlbCBpbiAxOTg2IOKAlCBzdWNoIGEgcmFwaWQgc2x1bXAgaGFzIG5ldmVyIGJlZW4gc2VlbiB3aGVuIGNvc3RzIHdlcmUgc28gaGlnaCBhbmQgd2l0aCBvaWwgcHJvZHVjdGlvbiBhbHJlYWR5IGluIGRlY2xpbmUuPC9wPg0KPHA+SVMgUFJPRFVDVElPTiBJTiBERUNMSU5FIEJFQ0FVU0UgTk9SVEggU0VBIElTIFJVTk5JTkcgRFJZPyBCVVQgRE9FUyBUSElTIERPV05UVVJOIFJJU0sgTUFTU0lWRSBBQ0NFTEVSQVRJT04gT0YgVEhBVCBERUNMSU5FIEJFQ0FVU0UgRVhQTE9SQVRJT04vREVWRUxPUE1FTlQgSVMgQkVJTkcgQ1VSVEFJTEVEPC9wPg0KPHA+T25lIGV4ZWN1dGl2ZSBmcm9tIGFuIG9pbCBtYWpvciBzYXlzOiDigJxXZSBoYXZlIG5ldmVyIGJlZm9yZSBleHBlcmllbmNlZCBhbnl0aGluZyBsaWtlIHRoaXMgYmVmb3JlIOKAlCBldmVyeXRoaW5nIGlzIGhhcHBlbmluZyBhdCB0aGUgc2FtZSB0aW1lLuKAnTwvcD4NCjxwPkFjY29yZGluZyB0byBhbmFseXNpcyBieSBFd2FuIE11cnJheSBvZiB0aGUgY29ycG9yYXRlIGhlYWx0aCBtb25pdG9yIENvbXBhbnkgV2F0Y2gsIG91dCBvZiB0aGUgMTI2IG9pbCBleHBsb3JhdGlvbiBhbmQgcHJvZHVjdGlvbiBjb21wYW5pZXMgbGlzdGVkIGluIExvbmRvbiwgNzcgcGVyIGNlbnQgYXJlIG5vdyBsb3NzLW1ha2luZy4gQVJFIEFMTCBUSEVTRSBJTiBUSEUgTk9SVEggU0VBPz8/IFRvdGFsIGxvc3NlcyBzdGFuZCBhdCDCozYuNmJuLjwvcD4NCjxwPjxhIGhyZWY9Imh0dHA6Ly9vaWxhbmRnYXN1ay5jby51ay9lY29ub21pYy1yZXBvcnQtMjAxNS5jZm0iIHRpdGxlPSJPaWwgJmFtcDsgR2FzIEVjb25vbWljIFJlcG9ydCAyMDE1Ij5GaWd1cmVzIHJlbGVhc2VkIGxhc3Qgd2VlayBieSBPaWwgJmFtcDsgR2FzIFVLPC9hPiwgdGhlIGluZHVzdHJ5IGJvZHksIHNob3cgdGhhdCBmb3IgdGhlIGZpcnN0IHRpbWUgc2luY2UgMTk3Nywgd2hlbiBOb3J0aCBTZWEgb2lsIHdhcyBzdGlsbCB5b3VuZywgaXQgbm93IGNvc3RzIG1vcmUgdG8gcHVtcCBvaWwgb3V0IG9mIHRoZSBzZWFiZWQgdGhhbiBpdCBnZW5lcmF0ZXMgaW4gcG9zdC10YXggcmV2ZW51ZT8/PyBQUk9GSVQ/Pz8uIDwvcD4NCjxwPjxhIGhyZWY9Imh0dHA6Ly9wdWJsaWMud29vZG1hYy5jb20vcHVibGljL21lZGlhLWNlbnRyZS8xMjUyOTI1NCIgdGl0bGU9IkxvdyBvaWwgcHJpY2UgYWNjZWxlcmF0aW5nIGRlY29tbWlzc2lvbmluZyBpbiB0aGUgVUtDUyAtIFdvb2RNYWMiPmFuYWx5c2lzIHB1Ymxpc2hlZCBsYXN0IHdlZWsgYnkgV29vZCBNYWNrZW56aWU8L2E+c3VnZ2VzdHMgMTQwIGZpZWxkcyBtYXkgY2xvc2UgaW4gdGhlIG5leHQgZml2ZSB5ZWFycywgZXZlbiBpZiB0aGUgb2lsIHByaWNlIHJldHVybnMgdG8gJDg1IGEgYmFycmVsLjwvcD4NCjxwPkRpZmZlcmVudCBjb21wYW5pZXMgb3BlcmF0aW5nIGluIHRoZSBOb3J0aCBTZWEgaGF2ZSByZXNwb25kZWQgdG8gdGhpcyBpbiBkaWZmZXJlbnQgd2F5cy48L3A+DQo8cD5PaWwgbWFqb3JzIHRlbmQgdG8gaGF2ZSBhbiBleGl0IHJvdXRlLiBUaGV5IGhhdmUgY3V0IGpvYnMg4oCUIDUsNTAwIGhhdmUgYmVlbiBsb3N0IHNvIGZhciwgd2l0aCBleGVjdXRpdmVzIHdhcm5pbmcgb2YgYW5vdGhlciAxMCwwMDAgaW4gdGhlIGNvbWluZyB5ZWFycyDigJQgYW5kIGFyZSBub3cgdHJ5aW5nIHRvIFJFRFVDRSBPVEhFUiBDT1NUUyBTVUNIIEFTPz8/IGN1dCBjb3N0cy4gQnV0IGlmIGFsbCBlbHNlIGZhaWxzLCB0aGV5IGNhbiBzaW1wbHkgZXhpdCB0aGUgcmVnaW9uIGFuZCBmb2N1cyB0aGVpciBlZmZvcnRzIG9uIGNoZWFwZXIgZmllbGRzIGVsc2V3aGVyZSBpbiB0aGUgd29ybGQuPC9wPg0KPHA+RnJhbmNl4oCZcyBUb3RhbCBsYXN0IG1vbnRoIGFncmVlZCB0byBzZWxsIE5vcnRoIFNlYSA8YSBocmVmPSJodHRwOi8vd3d3LmZ0LmNvbS9jbXMvcy8wLzg5MWQ2ZTU4LTRjMTktMTFlNS1iNTU4LThhOTcyMjk3NzE4OS5odG1sI2F4enozbG5Zckwzd2IiIHRpdGxlPSJGcmVuY2ggb2lsIG1ham9yIFRvdGFsIHRvIHNlbGwgJDkwMG0gb2YgTm9ydGggU2VhIGFzc2V0cyAtIEZUIj5nYXMgdGVybWluYWxzIGFuZCBwaXBlbGluZXM8L2E+IGluIGEgJDkwMG0gZGVhbCwgd2hpbGUgRW9uLCB0aGUgR2VybWFuIHV0aWxpdHkgY29tcGFueSwgaXMgbG9va2luZyBmb3IgYnV5ZXJzIGZvciBzb21lIG9mIGl0cyBhc3NldHMgaW4gdGhlIHJlZ2lvbi48L3A+DQo8cD5TbWFsbGVyIGNvbXBhbmllcyBob3dldmVyIHRlbmQgbm90IHRvIGhhdmUgdGhhdCBvcHRpb24uIEZvciBzb21lLCB0aGUgb2lsIHByaWNlIHBsdW5nZSBtZWFucyBzcXVlZXppbmcgY29zdHMgYXMgbXVjaCBhcyBwb3NzaWJsZSBpbiBhbiBlZmZvcnQgdG8gc3RheSBhZmxvYXQgdW50aWwgaXQgcmVib3VuZHMuPC9wPg0KPHA+RE8gQ0FQRVggQ1VUIEZJUlNUIEJZIEUgTlFVRVNUPC9wPg0KPHA+VEhFTiBFWFBMQUlOIEZPQ1VTIE9OIEVYVFJBQ1RJTkcgTU9SRSBGUk9NIEZFVyBGSUVMRFMgSVQgSVMgRk9DVVNJTkcgT04sIEFORCBDVVRUSU5HIE9QRVg8L3A+DQo8cD48L3A+DQo8cD5FbnF1ZXN0LCBmb3IgZXhhbXBsZSwgaGFzIG1hZGUgYSB2aXJ0dWUgZnJvbSBnZXR0aW5nIG1vcmUgb2lsIG91dCBvZiBtYXR1cmUgZmllbGRzIHRoYW4gbWFqb3JzIGNhbi4gSXRzIFRoaXN0bGUgb2lsIGZpZWxkIHdhcyBvbmNlIG93bmVkIGJ5IEJQLCBidXQgYXMgcHJvZHVjdGlvbiBkZWNsaW5lZCwgRW5xdWVzdCBzdGVwcGVkIGluLCBhbmQgYnkgbGFzdCB5ZWFyIHRoZSBjb21wYW55IG1hbmFnZWQgdG8gZXh0cmFjdCAzbSBiYXJyZWxzIGZyb20gaXQgZm9yIHRoZSBmaXJzdCB0aW1lIHNpbmNlIDE5OTcuPC9wPg0KPHA+T25lIHdheSB0byBkbyBzbz8/Pz8gaXMgdG8gY3V0IGNvc3RzIGFnZ3Jlc3NpdmVseS4gRW5xdWVzdCBoYXMgbG9va2VkIGF0IGV2ZXJ5dGhpbmcgaW5jbHVkaW5nIHJlbG9jYXRpbmcgaXRzIGhlbGljb3B0ZXIgc2VydmljZSBpbiBhbiBlZmZvcnQgdG8gc2F2ZSBtb25leS4gQW5kIG1vcmUgaXMgbGlrZWx5IHRvIGNvbWUuPC9wPg0KPHA+QW1qYWQgQnNlc2l1LCB0aGUgY29tcGFueeKAmXMgY2hpZWYgZXhlY3V0aXZlLCBzYWlkOiDigJxUaGUgaW5kdXN0cnkgd2VudCB0aHJvdWdoIGxhc3QgeWVhcuKAmXMgY3ljbGUgaW4gYSBsaXR0bGUgYml0IG9mIGEgdGltaWQgbWFubmVyLiBUaGVyZSB3YXMgYSBmZWVsaW5nIHRoYXQgbWF5YmUgdGhpbmdzIGNhbiBjb21lIGJhY2sgbGlrZSB0aGV5IGRpZCBpbiAyMDA4IHRvIDIwMDk/Pz8sIGJ1dCB0aGF0IGhhc27igJl0IGhhcHBlbmVkLuKAnTwvcD4NCjxwPkJ1dCBmb3IgbWFueSBzbWFsbGVyIGNvbXBhbmllcywgY3V0dGluZyBjb3N0cyBoYXMgYWxzbyBpbnZvbHZlZCBjdXR0aW5nIGV4cGxvcmF0aW9uIGFuZCBkZXZlbG9wbWVudC4gT3V0IG9mIHRoZSAxMiBmaWVsZHMgaXQgb3ducyBpbiB3aGljaCBvaWwgaGFzIGJlZW4gZGlzY292ZXJlZCwgRW5xdWVzdCBpcyBjdXJyZW50bHkgZGV2ZWxvcGluZyBvbmx5IHR3by4gVGhpcyBpcyBhIHRyZW5kIG92ZXIgdGhlIGluZHVzdHJ5IGFzIGEgd2hvbGU6IGluIDIwMDggdGhlIGluZHVzdHJ5IGRyaWxsZWQgNDAgZXhwbG9yYXRpb24gd2VsbHMgaW4gdGhlIE5vcnRoIFNlYS4gTGFzdCB5ZWFyLCB0aGF0IGZpZ3VyZSB3YXMganVzdCAxNC48L3A+DQo8cD48L3A+DQo8cD5CVVQgU09NRSBDT01QQU5JRVMgQVJFIFBST0NFRURJTkcgV0lUSCBCSUcgUFJPSkVDVFM6IFNVQ0ggQVMgUFJFTUlFUiBBVCBTT0xBTjwvcD4NCjxwPklUIE1FQU5TIFNPTUUgQ09NUEFOSUVTIEFSRSBHRU5FUkFUSU5HIElOU1VGRklDSUVOVCBDQVNIIFRPIENPVkVSIENBUEVYLCBTTyBSSVNJTkcgREVCVDwvcD4NCjxwPk1hbnkgc21hbGxlciBwbGF5ZXJzIDxzcGFuIGNoYW5uZWw9IiEiPmxpa2UgRW5xdWVzdCA8L3NwYW4+aGF2ZSBhbHNvIGNvcGVkIGJ5IHRha2luZyBvbiBpbmNyZWFzaW5nIGFtb3VudHMgb2YgZGVidC48L3A+DQo8cD5QcmVtaWVyIE9pbCwgZm9yIGV4YW1wbGUsIGhhZCAkNDBtIG9mIG5ldCBjYXNoIGluIHRoZSBmaXJzdCBxdWFydGVyIG9mIDIwMDkuIEl0cyBsYXN0IHJlc3VsdHMgc2hvd2VkIHRoYXQgdGhpcyBoYXMgbm93IGJlY29tZSAkMi4xYm4gd29ydGggb2YgbmV0IGRlYnQsIGFuZCBhbmFseXN0cyBzYXkgdGhpcyB3aWxsIGtlZXAgcmlzaW5nIHVudGlsIGF0IGxlYXN0IHRoZSBlbmQgb2YgdGhlIHllYXIuIE5FVCBERUJUIFRPIEVCSVREQSBNVUxUSVBMRTwvcD4NCjxwPk92ZXIgdGhlIHNhbWUgcGVyaW9kLCB0aGUgY29tcGFueeKAmXMgb2lsIHJlc2VydmVzIGhhdmUgZmFsbGVuIHNsaWdodGx5IGZyb20gMjI4bSBiYXJyZWxzIG9mIG9pbCBlcXVpdmFsZW50IHRvIDIyM20uIFRoZSBleHRyYSBkZWJ0IGhhcyBnb25lIHRvd2FyZHMgbWFpbnRhaW5pbmcgc3RvY2tzPz8/PyByYXRoZXIgdGhhbiBncm93aW5nIHRoZW0uPC9wPg0KPHA+V2hpbGUgbWFueSBvZiB0aGVzZSBzbWFsbGVyIHBsYXllcnMgYXJlIGhpZ2hseSBpbmRlYnRlZCwgc2V2ZXJhbCBhcmUgYWxzbyB3b3JraW5nIG9uIGNvbnNpZGVyYWJsZSBuZXcgZmluZHMsIHN1Y2ggYXMgUHJlbWllcuKAmXMgU29sYW4gZmllbGQgd2VzdCBvZiBTaGV0bGFuZC4gTW9zdCBjb21wYW5pZXMgaGF2ZSBtYW5hZ2VkIHRvIHJlbmVnb3RpYXRlIHRoZWlyIGRlYnQgaW4gcmVjZW50IG1vbnRocywgbWVhbmluZyBiYW5raW5nIGNvdmVuYW50cyBhcmUgdW5saWtlbHkgdG8gYmUgYnJlYWNoZWQgaW4gdGhlIGltbWVkaWF0ZSB0ZXJtLiA8L3A+DQo8cD5BbmQgd2l0aCB0YXggYnJlYWtzIHJlY2VudGx5IG9mZmVyZWQgYnkgdGhlIFRyZWFzdXJ5LCBhbnkgaW5jcmVhc2UgaW4gcmV2ZW51ZSBpcyBsaWtlbHkgdG8gbGVhZCB0byBhIHNpbWlsYXIgaW5jcmVhc2VzIGluIHByb2ZpdHMuPC9wPg0KPHA+TWFyayBXaWxzb24sIGFuIGFuYWx5c3QgYXQgSmVmZmVyaWVzIGJhbmssIHNhaWQ6IOKAnFJpc2luZyBuZXQgZGVidCBpcyB0YWtpbmcgdmFsdWUgYXdheSBmcm9tIGVxdWl0eSBob2xkZXJzIGJ1dCBzdG9ja3MgaGF2ZSBhbiB1cHNpZGUgaWYgdGhleSBjYW4gZGVsaXZlciB3aGF0IGlzIG9uIHRoZWlyIGJhbGFuY2Ugc2hlZXQuIFdIQVQgSVMgSEUgUkVGRVJSSU5HIFRPIE9OIEJBTEFOQ0UgU0hFRVQ/4oCdPC9wPg0KPHA+PC9wPg0KPHA+QnV0IGlmIHRoZSBuZXcgZmluZHMgZGlzYXBwb2ludCwgb3IgaWYgdGhlIG9pbCBwcmljZSByZW1haW5zIHRvbyBsb3cgdG8gY292ZXIgY29zdHMgd2hlbiBvaWwgc3RhcnRzIGNvbWluZyBvdXQgb2YgdGhlIGdyb3VuZCwgdGhleSBjb3VsZCBmaW5kIHRoZW1zZWx2ZXMgZmluYW5jaWFsbHkgZGlzdHJlc3NlZC4gSU4gV0hBVCBXQVkgV0FTIEZBSVJGSUVMRCBESVNUUkVTU0VEPzwvcD4NCjxwPk9uZSBjb21wYW55IHRoYXQgaGFzIGdvbmUgdGhyb3VnaCB0aGlzIHByb2Nlc3MgaXMgRmFpcmZpZWxkLCB3aGljaCBoYXMgbWFkZSB0aGUgZGVjaXNpb24gdG8gYWJhbmRvbiBpdHMgRHVubGluIEFscGhhIHBsYXRmb3JtIGFuZCB0dXJuIGl0c2VsZiBpbnRvIGEgZGVjb21taXNzaW9uaW5nIHNwZWNpYWxpc3QgaW5zdGVhZC48L3A+DQo8cCBjaGFubmVsPSIhIj5UaGlzIGNvdWxkIGJlIGEgc291bmQgYnVzaW5lc3MgZGVjaXNpb246PC9wPg0KPHA+V2hhdGV2ZXIgaGFwcGVucywgbm9ib2R5IGludm9sdmVkIGluIE5vcnRoIFNlYSBvaWwgdGhpbmtzIGl0IHdpbGwgbG9vayB0aGUgc2FtZSBieSB0aGUgZW5kIG9mIHRoZSBkZWNhZGUuIDwvcD4NCjxwPk9uZSBzZW5pb3IgZXhlY3V0aXZlIGZyb20gYSBtYWpvciBvaWwgY29tcGFueSBzYWlkOiDigJxMb3RzIG9mIGNvbXBhbmllcyBkbyBub3QgcmVhbGlzZSBpdCB5ZXQsIGJ1dCB0aGlzIGlzIHRoZSBiZWdpbm5pbmcgb2YgdGhlIGVuZC7igJ08L3A+DQo8L2JvZHk+DQo8L3RleHQ+DQo8L3N0b3J5Pg0KPC9kb2M+DQo="
const attributesMarkedDeletedTrue = "<ObjectMetadata><OutputChannels><DIFTcom><DIFTcomMarkDeleted>True</DIFTcomMarkDeleted></DIFTcom></OutputChannels></ObjectMetadata>"
const attributesMarkedDeletedFalse = "<ObjectMetadata><OutputChannels><DIFTcom><DIFTcomMarkDeleted>False</DIFTcomMarkDeleted></DIFTcom></OutputChannels></ObjectMetadata>"

const validListAttributes = "<!DOCTYPE ObjectMetadata SYSTEM \"/SysConfig/Classify/FTDWC2/classify.dtd\"><ObjectMetadata>	<FTcom>		<DIFTcomWebType>digitalList</DIFTcomWebType></FTcom><Lists>		<Title>Editor's pick</Title><LayoutHint>standard</LayoutHint></Lists></ObjectMetadata>"
const invalidListAttributes = "<!DOCTYPE ObjectMetadata SYSTEM \"/SysConfig/Classify/FTDWC2/classify.dtd\"><ObjectMetadata>	<FTcom>		<DIFTcomWebType>foobarList</DIFTcomWebType></FTcom><Lists>		<Title>Editor's pick</Title><LayoutHint>standard</LayoutHint></Lists></ObjectMetadata>"
const validFileTypeAttributes = "<?xml version=\"1.0\" encoding=\"UTF-8\"?><!DOCTYPE ObjectMetadata SYSTEM \"/SysConfig/Classify/FTStories/classify.dtd\"><ObjectMetadata><EditorialDisplayIndexing><DIBylineCopy>Kiran Stacey, Energy Correspondent</DIBylineCopy></EditorialDisplayIndexing><OutputChannels><DIFTcom><DIFTcomWebType>story</DIFTcomWebType></DIFTcom></OutputChannels><EditorialNotes><Sources><Source><SourceCode>FT</SourceCode></Source></Sources><Language>English</Language><Author>staceyk</Author><ObjectLocation>/Users/staceyk/North Sea companies analysis CO 15.xml</ObjectLocation></EditorialNotes><DataFactoryIndexing><ADRIS_MetaData><IndexSuccess>yes</IndexSuccess><StartTime>Wed Oct 21 10:11:53 GMT 2015</StartTime><EndTime>Wed Oct 21 10:11:57 GMT 2015</EndTime></ADRIS_MetaData></DataFactoryIndexing></ObjectMetadata>"
const invalidFileTypeAttributes = "<?xml version=\"1.0\" encoding=\"UTF-8\"?><!DOCTYPE ObjectMetadata SYSTEM \"/SysConfig/Classify/FTStories/classify.dtd\"><ObjectMetadata><EditorialDisplayIndexing><DIBylineCopy>Kiran Stacey, Energy Correspondent</DIBylineCopy></EditorialDisplayIndexing><OutputChannels><DIFTcom><DIFTcomWebType>story</DIFTcomWebType></DIFTcom></OutputChannels><EditorialNotes><Language>English</Language><Author>staceyk</Author><ObjectLocation>/Users/staceyk/North Sea companies analysis CO 15.PDF</ObjectLocation></EditorialNotes><DataFactoryIndexing><ADRIS_MetaData><IndexSuccess>yes</IndexSuccess><StartTime>Wed Oct 21 10:11:53 GMT 2015</StartTime><EndTime>Wed Oct 21 10:11:57 GMT 2015</EndTime></ADRIS_MetaData></DataFactoryIndexing></ObjectMetadata>"

var supportedSourceCodeAttributes = "<?xml version=\"1.0\" encoding=\"UTF-8\"?><!DOCTYPE ObjectMetadata SYSTEM \"/SysConfig/Classify/FTStories/classify.dtd\"><ObjectMetadata><EditorialDisplayIndexing><DILeadCompanies /><DITemporaryCompanies><DITemporaryCompany><DICoTempCode /><DICoTempDescriptor /><DICoTickerCode /></DITemporaryCompany></DITemporaryCompanies><DIFTSEGlobalClassifications /><DIStockExchangeIndices /><DIHotTopics /><DIHeadlineCopy>Global oil inventory stands at record level</DIHeadlineCopy><DIBylineCopy>Anjli Raval, Oil and Gas Correspondent</DIBylineCopy><DIFTNPSections /><DIFirstParCopy>International Energy Agency says crude stockpiles near 3bn barrels despite robust demand growth</DIFirstParCopy><DIMasterImgFileRef>/FT/Graphics/Online/Master_2048x1152/2015/08/New%20story%20of%20MAS_crude-oil_02.jpg?uuid=c61fcc32-61cd-11e5-9846-de406ccb37f2</DIMasterImgFileRef></EditorialDisplayIndexing><OutputChannels><DIFTN><DIFTNPublicationDate /><DIFTNZoneEdition /><DIFTNPage /><DIFTNTimeEdition /><DIFTNFronts /></DIFTN><DIFTcom><DIFTcomWebType>story</DIFTcomWebType><DIFTcomDisplayCodes><DIFTcomDisplayCodeRank1><DIFTcomDisplayCode title=\"Markets - Commodities\"><DIFTcomDisplayCodeFTCode>MKCO</DIFTcomDisplayCodeFTCode><DIFTcomDisplayCodeDescriptor>Markets - Commodities</DIFTcomDisplayCodeDescriptor><DIFTcomDisplayCodeNewsInDepth>True</DIFTcomDisplayCodeNewsInDepth><DIFTcomDisplayCodeSite>FTcom</DIFTcomDisplayCodeSite><DIFTcomDisplayCodeArticleType>News</DIFTcomDisplayCodeArticleType><DIFTcomDisplayCodeArticleBrand /><DIFTcomDisplayCodeEditorsTag>Commodities</DIFTcomDisplayCodeEditorsTag></DIFTcomDisplayCode></DIFTcomDisplayCodeRank1><DIFTcomDisplayCodeRank2><DIFTcomDisplayCode title=\"Markets - European Equities\"><DIFTcomDisplayCodeFTCode>MKEU</DIFTcomDisplayCodeFTCode><DIFTcomDisplayCodeDescriptor>Markets - European Equities</DIFTcomDisplayCodeDescriptor><DIFTcomDisplayCodeNewsInDepth>True</DIFTcomDisplayCodeNewsInDepth><DIFTcomDisplayCodeSite>FTcom</DIFTcomDisplayCodeSite><DIFTcomDisplayCodeArticleType>News</DIFTcomDisplayCodeArticleType><DIFTcomDisplayCodeArticleBrand /><DIFTcomDisplayCodeEditorsTag>European Equities</DIFTcomDisplayCodeEditorsTag></DIFTcomDisplayCode><DIFTcomDisplayCode title=\"Energy - Oil &amp; Gas\"><DIFTcomDisplayCodeFTCode>OG8E</DIFTcomDisplayCodeFTCode><DIFTcomDisplayCodeDescriptor>Energy - Oil &amp; Gas</DIFTcomDisplayCodeDescriptor><DIFTcomDisplayCodeNewsInDepth>True</DIFTcomDisplayCodeNewsInDepth><DIFTcomDisplayCodeSite>FTcom</DIFTcomDisplayCodeSite><DIFTcomDisplayCodeArticleType>News</DIFTcomDisplayCodeArticleType><DIFTcomDisplayCodeArticleBrand /><DIFTcomDisplayCodeEditorsTag>Oil &amp; Gas</DIFTcomDisplayCodeEditorsTag></DIFTcomDisplayCode><DIFTcomDisplayCode title=\"Energy\"><DIFTcomDisplayCodeFTCode>NDEM</DIFTcomDisplayCodeFTCode><DIFTcomDisplayCodeDescriptor>Energy</DIFTcomDisplayCodeDescriptor><DIFTcomDisplayCodeNewsInDepth>True</DIFTcomDisplayCodeNewsInDepth><DIFTcomDisplayCodeSite>FTcom</DIFTcomDisplayCodeSite><DIFTcomDisplayCodeArticleType>News</DIFTcomDisplayCodeArticleType><DIFTcomDisplayCodeArticleBrand /><DIFTcomDisplayCodeEditorsTag>Energy</DIFTcomDisplayCodeEditorsTag></DIFTcomDisplayCode><DIFTcomDisplayCode title=\"Companies\"><DIFTcomDisplayCodeFTCode>BNIP</DIFTcomDisplayCodeFTCode><DIFTcomDisplayCodeDescriptor>Companies</DIFTcomDisplayCodeDescriptor><DIFTcomDisplayCodeNewsInDepth>True</DIFTcomDisplayCodeNewsInDepth><DIFTcomDisplayCodeSite>FTcom</DIFTcomDisplayCodeSite><DIFTcomDisplayCodeArticleType>News</DIFTcomDisplayCodeArticleType><DIFTcomDisplayCodeArticleBrand /><DIFTcomDisplayCodeEditorsTag /></DIFTcomDisplayCode><DIFTcomDisplayCode title=\"Markets\"><DIFTcomDisplayCodeFTCode>MKIP</DIFTcomDisplayCodeFTCode><DIFTcomDisplayCodeDescriptor>Markets</DIFTcomDisplayCodeDescriptor><DIFTcomDisplayCodeNewsInDepth>True</DIFTcomDisplayCodeNewsInDepth><DIFTcomDisplayCodeSite>FTcom</DIFTcomDisplayCodeSite><DIFTcomDisplayCodeArticleType>News</DIFTcomDisplayCodeArticleType><DIFTcomDisplayCodeArticleBrand /><DIFTcomDisplayCodeEditorsTag /></DIFTcomDisplayCode><DIFTcomDisplayCode title=\"Markets - Equities Main\"><DIFTcomDisplayCodeFTCode>MKEQ</DIFTcomDisplayCodeFTCode><DIFTcomDisplayCodeDescriptor>Markets - Equities Main</DIFTcomDisplayCodeDescriptor><DIFTcomDisplayCodeNewsInDepth>True</DIFTcomDisplayCodeNewsInDepth><DIFTcomDisplayCodeSite>FTcom</DIFTcomDisplayCodeSite><DIFTcomDisplayCodeArticleType>News</DIFTcomDisplayCodeArticleType><DIFTcomDisplayCodeArticleBrand /><DIFTcomDisplayCodeEditorsTag>Equities</DIFTcomDisplayCodeEditorsTag></DIFTcomDisplayCode></DIFTcomDisplayCodeRank2></DIFTcomDisplayCodes><DIFTcomSubscriptionLevel>0</DIFTcomSubscriptionLevel><DIFTcomUpdateTimeStamp>False</DIFTcomUpdateTimeStamp><DIFTcomIndexAndSynd>false</DIFTcomIndexAndSynd><DIFTcomSafeToSyndicate>True</DIFTcomSafeToSyndicate><DIFTcomInitialPublication>20151113105031</DIFTcomInitialPublication><DIFTcomLastPublication>20151113132953</DIFTcomLastPublication><DIFTcomSuppresInlineAds>False</DIFTcomSuppresInlineAds><DIFTcomMap>True</DIFTcomMap><DIFTcomDisplayStyle>Normal</DIFTcomDisplayStyle><DIFTcomFeatureType>Normal</DIFTcomFeatureType><DIFTcomMarkDeleted>False</DIFTcomMarkDeleted><DIFTcomMakeUnlinkable>False</DIFTcomMakeUnlinkable><isBestStory>0</isBestStory><DIFTcomCMRId>3040370</DIFTcomCMRId><DIFTcomCMRHint /><DIFTcomCMR><DIFTcomCMRPrimarySection>Commodities</DIFTcomCMRPrimarySection><DIFTcomCMRPrimarySectionId>MTA1-U2VjdGlvbnM=</DIFTcomCMRPrimarySectionId><DIFTcomCMRPrimaryTheme>Oil</DIFTcomCMRPrimaryTheme><DIFTcomCMRPrimaryThemeId>ZmFmYTUxOTItMGZjZC00YmJkLWJlZTQtMmY3ZDZiOWZkYmYw-VG9waWNz</DIFTcomCMRPrimaryThemeId><DIFTcomCMRBrand /><DIFTcomCMRBrandId /><DIFTcomCMRGenre>News</DIFTcomCMRGenre><DIFTcomCMRGenreId>Nw==-R2VucmVz</DIFTcomCMRGenreId><DIFTcomCMRMediaType>Text</DIFTcomCMRMediaType><DIFTcomCMRMediaTypeId>ZjMwY2E2NjctMDA1Ni00ZTk4LWI0MWUtZjk5MTk2ZTMyNGVm-TWVkaWFUeXBlcw==</DIFTcomCMRMediaTypeId></DIFTcomCMR><DIFTcomECPositionInText>Default</DIFTcomECPositionInText><DIFTcomHideECLevel1>False</DIFTcomHideECLevel1><DIFTcomHideECLevel2>False</DIFTcomHideECLevel2><DIFTcomHideECLevel3>False</DIFTcomHideECLevel3><DIFTcomDiscussion>True</DIFTcomDiscussion><DIFTcomArticleImage>Article size</DIFTcomArticleImage></DIFTcom><DISyndication><DISyndBeenCopied>False</DISyndBeenCopied><DISyndEdition>USA</DISyndEdition><DISyndStar>01</DISyndStar><DISyndChannel /><DISyndArea /><DISyndCategory /></DISyndication></OutputChannels><EditorialNotes><Language>English</Language><Author>ravala</Author><Guides /><Editor /><Sources><Source title=\"Financial Times\"><SourceCode>FT</SourceCode><SourceDescriptor>Financial Times</SourceDescriptor><SourceOnlineInclusion>True</SourceOnlineInclusion><SourceCanBeSyndicated>True</SourceCanBeSyndicated></Source></Sources><WordCount>670</WordCount><CreationDate>20151113095342</CreationDate><EmbargoDate /><ExpiryDate>20151113095342</ExpiryDate><ObjectLocation>/FT/Content/Markets/Stories/Live/IEA MA 13.xml</ObjectLocation><OriginatingStory /><CCMS><CCMSCommissionRefNo /><CCMSContributorRefNo>CB-0000844</CCMSContributorRefNo><CCMSContributorFullName>Anjli Raval</CCMSContributorFullName><CCMSContributorInclude /><CCMSContributorRights>3</CCMSContributorRights><CCMSFilingDate /><CCMSProposedPublishingDate /></CCMS></EditorialNotes><WiresIndexing><category /><Keyword /><char_count /><priority /><basket /><title /><Version /><story_num /><file_name /><serviceid /><entry_date /><ref_field /><take_num /></WiresIndexing><DataFactoryIndexing><ADRIS_MetaData><IndexSuccess>yes</IndexSuccess><StartTime>Fri Nov 13 13:29:53 GMT 2015</StartTime><EndTime>Fri Nov 13 13:29:56 GMT 2015</EndTime></ADRIS_MetaData><DFMajorCompanies><DFMajorCompany><DFCoMajScore>100</DFCoMajScore><DFCoMajDescriptor>International_Energy_Agency</DFCoMajDescriptor><DFCoMajFTCode>IEIEA00000</DFCoMajFTCode><Version>1</Version><DFCoMajTickerSymbol /><DFCoMajTickerExchangeCountry /><DFCoMajTickerExchangeCode /><DFCoMajFTMWTickercode /><DFCoMajSEDOL /><DFCoMajISIN /><DFCoMajCOFlag>O</DFCoMajCOFlag></DFMajorCompany></DFMajorCompanies><DFMinorCompanies /><DFNAICS><DFNAIC><DFNAICSFTCode>N52</DFNAICSFTCode><DFNAICSDescriptor>Finance_&amp;_Insurance</DFNAICSDescriptor><Version>1997</Version></DFNAIC><DFNAIC><DFNAICSFTCode>N523</DFNAICSFTCode><DFNAICSDescriptor>Security_Commodity_Contracts_&amp;_Like_Activity</DFNAICSDescriptor><Version>1997</Version></DFNAIC><DFNAIC><DFNAICSFTCode>N52321</DFNAICSFTCode><DFNAICSDescriptor>Securities_&amp;_Commodity_Exchanges</DFNAICSDescriptor><Version>1997</Version></DFNAIC><DFNAIC><DFNAICSFTCode>N9261</DFNAICSFTCode><DFNAICSDescriptor>Admin_of_Economic_Programs</DFNAICSDescriptor><Version>1997</Version></DFNAIC><DFNAIC><DFNAICSFTCode>N92613</DFNAICSFTCode><DFNAICSDescriptor>Regulation_&amp;_Admin_of_Utilities</DFNAICSDescriptor><Version>1997</Version></DFNAIC></DFNAICS><DFWPMIndustries /><DFFTSEGlobalClassifications /><DFStockExchangeIndices /><DFSubjects /><DFCountries /><DFRegions /><DFWPMRegions /><DFProvinces /><DFFTcomDisplayCodes /><DFFTSections /><DFWebRegions /></DataFactoryIndexing></ObjectMetadata>"
var supportedSourceCodeAttributesContentPlaceholder = "<?xml version=\"1.0\" encoding=\"UTF-8\"?><!DOCTYPE ObjectMetadata SYSTEM \"/SysConfig/Classify/FTStories/classify.dtd\"><ObjectMetadata><EditorialDisplayIndexing><DILeadCompanies /><DITemporaryCompanies><DITemporaryCompany><DICoTempCode /><DICoTempDescriptor /><DICoTickerCode /></DITemporaryCompany></DITemporaryCompanies><DIFTSEGlobalClassifications /><DIStockExchangeIndices /><DIHotTopics /><DIHeadlineCopy>Global oil inventory stands at record level</DIHeadlineCopy><DIBylineCopy>Anjli Raval, Oil and Gas Correspondent</DIBylineCopy><DIFTNPSections /><DIFirstParCopy>International Energy Agency says crude stockpiles near 3bn barrels despite robust demand growth</DIFirstParCopy><DIMasterImgFileRef>/FT/Graphics/Online/Master_2048x1152/2015/08/New%20story%20of%20MAS_crude-oil_02.jpg?uuid=c61fcc32-61cd-11e5-9846-de406ccb37f2</DIMasterImgFileRef></EditorialDisplayIndexing><OutputChannels><DIFTN><DIFTNPublicationDate /><DIFTNZoneEdition /><DIFTNPage /><DIFTNTimeEdition /><DIFTNFronts /></DIFTN><DIFTcom><DIFTcomWebType>story</DIFTcomWebType><DIFTcomDisplayCodes><DIFTcomDisplayCodeRank1><DIFTcomDisplayCode title=\"Markets - Commodities\"><DIFTcomDisplayCodeFTCode>MKCO</DIFTcomDisplayCodeFTCode><DIFTcomDisplayCodeDescriptor>Markets - Commodities</DIFTcomDisplayCodeDescriptor><DIFTcomDisplayCodeNewsInDepth>True</DIFTcomDisplayCodeNewsInDepth><DIFTcomDisplayCodeSite>FTcom</DIFTcomDisplayCodeSite><DIFTcomDisplayCodeArticleType>News</DIFTcomDisplayCodeArticleType><DIFTcomDisplayCodeArticleBrand /><DIFTcomDisplayCodeEditorsTag>Commodities</DIFTcomDisplayCodeEditorsTag></DIFTcomDisplayCode></DIFTcomDisplayCodeRank1><DIFTcomDisplayCodeRank2><DIFTcomDisplayCode title=\"Markets - European Equities\"><DIFTcomDisplayCodeFTCode>MKEU</DIFTcomDisplayCodeFTCode><DIFTcomDisplayCodeDescriptor>Markets - European Equities</DIFTcomDisplayCodeDescriptor><DIFTcomDisplayCodeNewsInDepth>True</DIFTcomDisplayCodeNewsInDepth><DIFTcomDisplayCodeSite>FTcom</DIFTcomDisplayCodeSite><DIFTcomDisplayCodeArticleType>News</DIFTcomDisplayCodeArticleType><DIFTcomDisplayCodeArticleBrand /><DIFTcomDisplayCodeEditorsTag>European Equities</DIFTcomDisplayCodeEditorsTag></DIFTcomDisplayCode><DIFTcomDisplayCode title=\"Energy - Oil &amp; Gas\"><DIFTcomDisplayCodeFTCode>OG8E</DIFTcomDisplayCodeFTCode><DIFTcomDisplayCodeDescriptor>Energy - Oil &amp; Gas</DIFTcomDisplayCodeDescriptor><DIFTcomDisplayCodeNewsInDepth>True</DIFTcomDisplayCodeNewsInDepth><DIFTcomDisplayCodeSite>FTcom</DIFTcomDisplayCodeSite><DIFTcomDisplayCodeArticleType>News</DIFTcomDisplayCodeArticleType><DIFTcomDisplayCodeArticleBrand /><DIFTcomDisplayCodeEditorsTag>Oil &amp; Gas</DIFTcomDisplayCodeEditorsTag></DIFTcomDisplayCode><DIFTcomDisplayCode title=\"Energy\"><DIFTcomDisplayCodeFTCode>NDEM</DIFTcomDisplayCodeFTCode><DIFTcomDisplayCodeDescriptor>Energy</DIFTcomDisplayCodeDescriptor><DIFTcomDisplayCodeNewsInDepth>True</DIFTcomDisplayCodeNewsInDepth><DIFTcomDisplayCodeSite>FTcom</DIFTcomDisplayCodeSite><DIFTcomDisplayCodeArticleType>News</DIFTcomDisplayCodeArticleType><DIFTcomDisplayCodeArticleBrand /><DIFTcomDisplayCodeEditorsTag>Energy</DIFTcomDisplayCodeEditorsTag></DIFTcomDisplayCode><DIFTcomDisplayCode title=\"Companies\"><DIFTcomDisplayCodeFTCode>BNIP</DIFTcomDisplayCodeFTCode><DIFTcomDisplayCodeDescriptor>Companies</DIFTcomDisplayCodeDescriptor><DIFTcomDisplayCodeNewsInDepth>True</DIFTcomDisplayCodeNewsInDepth><DIFTcomDisplayCodeSite>FTcom</DIFTcomDisplayCodeSite><DIFTcomDisplayCodeArticleType>News</DIFTcomDisplayCodeArticleType><DIFTcomDisplayCodeArticleBrand /><DIFTcomDisplayCodeEditorsTag /></DIFTcomDisplayCode><DIFTcomDisplayCode title=\"Markets\"><DIFTcomDisplayCodeFTCode>MKIP</DIFTcomDisplayCodeFTCode><DIFTcomDisplayCodeDescriptor>Markets</DIFTcomDisplayCodeDescriptor><DIFTcomDisplayCodeNewsInDepth>True</DIFTcomDisplayCodeNewsInDepth><DIFTcomDisplayCodeSite>FTcom</DIFTcomDisplayCodeSite><DIFTcomDisplayCodeArticleType>News</DIFTcomDisplayCodeArticleType><DIFTcomDisplayCodeArticleBrand /><DIFTcomDisplayCodeEditorsTag /></DIFTcomDisplayCode><DIFTcomDisplayCode title=\"Markets - Equities Main\"><DIFTcomDisplayCodeFTCode>MKEQ</DIFTcomDisplayCodeFTCode><DIFTcomDisplayCodeDescriptor>Markets - Equities Main</DIFTcomDisplayCodeDescriptor><DIFTcomDisplayCodeNewsInDepth>True</DIFTcomDisplayCodeNewsInDepth><DIFTcomDisplayCodeSite>FTcom</DIFTcomDisplayCodeSite><DIFTcomDisplayCodeArticleType>News</DIFTcomDisplayCodeArticleType><DIFTcomDisplayCodeArticleBrand /><DIFTcomDisplayCodeEditorsTag>Equities</DIFTcomDisplayCodeEditorsTag></DIFTcomDisplayCode></DIFTcomDisplayCodeRank2></DIFTcomDisplayCodes><DIFTcomSubscriptionLevel>0</DIFTcomSubscriptionLevel><DIFTcomUpdateTimeStamp>False</DIFTcomUpdateTimeStamp><DIFTcomIndexAndSynd>false</DIFTcomIndexAndSynd><DIFTcomSafeToSyndicate>True</DIFTcomSafeToSyndicate><DIFTcomInitialPublication>20151113105031</DIFTcomInitialPublication><DIFTcomLastPublication>20151113132953</DIFTcomLastPublication><DIFTcomSuppresInlineAds>False</DIFTcomSuppresInlineAds><DIFTcomMap>True</DIFTcomMap><DIFTcomDisplayStyle>Normal</DIFTcomDisplayStyle><DIFTcomFeatureType>Normal</DIFTcomFeatureType><DIFTcomMarkDeleted>False</DIFTcomMarkDeleted><DIFTcomMakeUnlinkable>False</DIFTcomMakeUnlinkable><isBestStory>0</isBestStory><DIFTcomCMRId>3040370</DIFTcomCMRId><DIFTcomCMRHint /><DIFTcomCMR><DIFTcomCMRPrimarySection>Commodities</DIFTcomCMRPrimarySection><DIFTcomCMRPrimarySectionId>MTA1-U2VjdGlvbnM=</DIFTcomCMRPrimarySectionId><DIFTcomCMRPrimaryTheme>Oil</DIFTcomCMRPrimaryTheme><DIFTcomCMRPrimaryThemeId>ZmFmYTUxOTItMGZjZC00YmJkLWJlZTQtMmY3ZDZiOWZkYmYw-VG9waWNz</DIFTcomCMRPrimaryThemeId><DIFTcomCMRBrand /><DIFTcomCMRBrandId /><DIFTcomCMRGenre>News</DIFTcomCMRGenre><DIFTcomCMRGenreId>Nw==-R2VucmVz</DIFTcomCMRGenreId><DIFTcomCMRMediaType>Text</DIFTcomCMRMediaType><DIFTcomCMRMediaTypeId>ZjMwY2E2NjctMDA1Ni00ZTk4LWI0MWUtZjk5MTk2ZTMyNGVm-TWVkaWFUeXBlcw==</DIFTcomCMRMediaTypeId></DIFTcomCMR><DIFTcomECPositionInText>Default</DIFTcomECPositionInText><DIFTcomHideECLevel1>False</DIFTcomHideECLevel1><DIFTcomHideECLevel2>False</DIFTcomHideECLevel2><DIFTcomHideECLevel3>False</DIFTcomHideECLevel3><DIFTcomDiscussion>True</DIFTcomDiscussion><DIFTcomArticleImage>Article size</DIFTcomArticleImage></DIFTcom><DISyndication><DISyndBeenCopied>False</DISyndBeenCopied><DISyndEdition>USA</DISyndEdition><DISyndStar>01</DISyndStar><DISyndChannel /><DISyndArea /><DISyndCategory /></DISyndication></OutputChannels><EditorialNotes><Language>English</Language><Author>ravala</Author><Guides /><Editor /><Sources><Source title=\"Financial Times\"><SourceCode>ContentPlaceholder</SourceCode><SourceDescriptor>Financial Times</SourceDescriptor><SourceOnlineInclusion>True</SourceOnlineInclusion><SourceCanBeSyndicated>True</SourceCanBeSyndicated></Source></Sources><WordCount>670</WordCount><CreationDate>20151113095342</CreationDate><EmbargoDate /><ExpiryDate>20151113095342</ExpiryDate><ObjectLocation>/FT/Content/Markets/Stories/Live/IEA MA 13.xml</ObjectLocation><OriginatingStory /><CCMS><CCMSCommissionRefNo /><CCMSContributorRefNo>CB-0000844</CCMSContributorRefNo><CCMSContributorFullName>Anjli Raval</CCMSContributorFullName><CCMSContributorInclude /><CCMSContributorRights>3</CCMSContributorRights><CCMSFilingDate /><CCMSProposedPublishingDate /></CCMS></EditorialNotes><WiresIndexing><category /><Keyword /><char_count /><priority /><basket /><title /><Version /><story_num /><file_name /><serviceid /><entry_date /><ref_field /><take_num /></WiresIndexing><DataFactoryIndexing><ADRIS_MetaData><IndexSuccess>yes</IndexSuccess><StartTime>Fri Nov 13 13:29:53 GMT 2015</StartTime><EndTime>Fri Nov 13 13:29:56 GMT 2015</EndTime></ADRIS_MetaData><DFMajorCompanies><DFMajorCompany><DFCoMajScore>100</DFCoMajScore><DFCoMajDescriptor>International_Energy_Agency</DFCoMajDescriptor><DFCoMajFTCode>IEIEA00000</DFCoMajFTCode><Version>1</Version><DFCoMajTickerSymbol /><DFCoMajTickerExchangeCountry /><DFCoMajTickerExchangeCode /><DFCoMajFTMWTickercode /><DFCoMajSEDOL /><DFCoMajISIN /><DFCoMajCOFlag>O</DFCoMajCOFlag></DFMajorCompany></DFMajorCompanies><DFMinorCompanies /><DFNAICS><DFNAIC><DFNAICSFTCode>N52</DFNAICSFTCode><DFNAICSDescriptor>Finance_&amp;_Insurance</DFNAICSDescriptor><Version>1997</Version></DFNAIC><DFNAIC><DFNAICSFTCode>N523</DFNAICSFTCode><DFNAICSDescriptor>Security_Commodity_Contracts_&amp;_Like_Activity</DFNAICSDescriptor><Version>1997</Version></DFNAIC><DFNAIC><DFNAICSFTCode>N52321</DFNAICSFTCode><DFNAICSDescriptor>Securities_&amp;_Commodity_Exchanges</DFNAICSDescriptor><Version>1997</Version></DFNAIC><DFNAIC><DFNAICSFTCode>N9261</DFNAICSFTCode><DFNAICSDescriptor>Admin_of_Economic_Programs</DFNAICSDescriptor><Version>1997</Version></DFNAIC><DFNAIC><DFNAICSFTCode>N92613</DFNAICSFTCode><DFNAICSDescriptor>Regulation_&amp;_Admin_of_Utilities</DFNAICSDescriptor><Version>1997</Version></DFNAIC></DFNAICS><DFWPMIndustries /><DFFTSEGlobalClassifications /><DFStockExchangeIndices /><DFSubjects /><DFCountries /><DFRegions /><DFWPMRegions /><DFProvinces /><DFFTcomDisplayCodes /><DFFTSections /><DFWebRegions /></DataFactoryIndexing></ObjectMetadata>"
var unsupportedSourceCodeAttributes = "<?xml version=\"1.0\" encoding=\"UTF-8\"?><!DOCTYPE ObjectMetadata SYSTEM \"/SysConfig/Classify/FTStories/classify.dtd\"><ObjectMetadata><EditorialDisplayIndexing><DILeadCompanies /><DITemporaryCompanies><DITemporaryCompany><DICoTempCode /><DICoTempDescriptor /><DICoTickerCode /></DITemporaryCompany></DITemporaryCompanies><DIFTSEGlobalClassifications /><DIStockExchangeIndices /><DIHotTopics /><DIHeadlineCopy /><DIBylineCopy>Pilita Clark and Aimee Keane</DIBylineCopy><DIFTNPSections /><DIFirstParCopy /><DIMasterImgFileRef>/FT/Graphics/Online/Master_2048x1152/2015/01/MAS_m238260_20150118210430.jpg?uuid=987d1f18-9f55-11e4-a849-00144feab7de</DIMasterImgFileRef></EditorialDisplayIndexing><OutputChannels><DIFTN><DIFTNPublicationDate /><DIFTNZoneEdition /><DIFTNPage /><DIFTNTimeEdition /><DIFTNFronts /></DIFTN><DIFTcom><DIFTcomWebType>story</DIFTcomWebType><DIFTcomDisplayCodes><DIFTcomDisplayCodeRank1 /><DIFTcomDisplayCodeRank2 /></DIFTcomDisplayCodes><DIFTcomSubscriptionLevel>0</DIFTcomSubscriptionLevel><DIFTcomUpdateTimeStamp>False</DIFTcomUpdateTimeStamp><DIFTcomIndexAndSynd>true</DIFTcomIndexAndSynd><DIFTcomSafeToSyndicate>True</DIFTcomSafeToSyndicate><DIFTcomInitialPublication>20150306111836</DIFTcomInitialPublication><DIFTcomLastPublication>20151112154414</DIFTcomLastPublication><DIFTcomSuppresInlineAds>False</DIFTcomSuppresInlineAds><DIFTcomMap>True</DIFTcomMap><DIFTcomDisplayStyle>Normal</DIFTcomDisplayStyle><DIFTcomFeatureType>Normal</DIFTcomFeatureType><DIFTcomMarkDeleted>False</DIFTcomMarkDeleted><DIFTcomMakeUnlinkable>False</DIFTcomMakeUnlinkable><isBestStory>0</isBestStory><DIFTcomCMRId>2884440</DIFTcomCMRId><DIFTcomCMRHint /><DIFTcomCMR><DIFTcomCMRPrimarySection>Timelines</DIFTcomCMRPrimarySection><DIFTcomCMRPrimarySectionId>NzBlYjllYjgtZjlmZC00YmQ5LTg2MzItOTU2ZjFhM2RjNWFm-U2VjdGlvbnM=</DIFTcomCMRPrimarySectionId><DIFTcomCMRPrimaryTheme>Climate change</DIFTcomCMRPrimaryTheme><DIFTcomCMRPrimaryThemeId>Ng==-VG9waWNz</DIFTcomCMRPrimaryThemeId><DIFTcomCMRBrand /><DIFTcomCMRBrandId /><DIFTcomCMRGenre>News</DIFTcomCMRGenre><DIFTcomCMRGenreId>Nw==-R2VucmVz</DIFTcomCMRGenreId><DIFTcomCMRMediaType>Text</DIFTcomCMRMediaType><DIFTcomCMRMediaTypeId>ZjMwY2E2NjctMDA1Ni00ZTk4LWI0MWUtZjk5MTk2ZTMyNGVm-TWVkaWFUeXBlcw==</DIFTcomCMRMediaTypeId></DIFTcomCMR><DIFTcomECPositionInText>Default</DIFTcomECPositionInText><DIFTcomHideECLevel1>True</DIFTcomHideECLevel1><DIFTcomHideECLevel2>True</DIFTcomHideECLevel2><DIFTcomHideECLevel3>True</DIFTcomHideECLevel3><DIFTcomDiscussion>True</DIFTcomDiscussion><DIFTcomArticleImage>Primary size</DIFTcomArticleImage></DIFTcom><DISyndication><DISyndBeenCopied>False</DISyndBeenCopied><DISyndEdition>USA</DISyndEdition><DISyndStar>01</DISyndStar><DISyndChannel /><DISyndArea /><DISyndCategory /></DISyndication></OutputChannels><EditorialNotes><Language>English</Language><Author>kwongr</Author><Guides /><Editor /><Sources><Source title=\"FT Sitebuild\"><SourceCode>FTSB</SourceCode><SourceDescriptor>FT Sitebuild</SourceDescriptor><SourceOnlineInclusion>True</SourceOnlineInclusion><SourceCanBeSyndicated>False</SourceCanBeSyndicated></Source></Sources><WordCount>38</WordCount><CreationDate>20150306111437</CreationDate><EmbargoDate /><ExpiryDate>20150306111437</ExpiryDate><ObjectLocation>/FT/Content/Interactive/Stories/Live/Paris timeline.xml</ObjectLocation><OriginatingStory /><CCMS><CCMSCommissionRefNo /><CCMSContributorRefNo>CB-0000698</CCMSContributorRefNo><CCMSContributorFullName>Pilita Clark</CCMSContributorFullName><CCMSContributorInclude /><CCMSContributorRights>3</CCMSContributorRights><CCMSFilingDate /><CCMSProposedPublishingDate /></CCMS></EditorialNotes><WiresIndexing><category /><Keyword /><char_count /><priority /><basket /><title /><Version /><story_num /><file_name /><serviceid /><entry_date /><ref_field /><take_num /></WiresIndexing><DataFactoryIndexing><ADRIS_MetaData><IndexSuccess>yes</IndexSuccess><StartTime>Thu Nov 12 15:44:15 GMT 2015</StartTime><EndTime>Thu Nov 12 15:44:15 GMT 2015</EndTime></ADRIS_MetaData><DFMajorCompanies /><DFMinorCompanies /><DFNAICS /><DFWPMIndustries /><DFFTSEGlobalClassifications /><DFStockExchangeIndices /><DFSubjects><DFSubject><DFSUFTCode>ON15</DFSUFTCode><DFSUDescriptor>Environment</DFSUDescriptor><Version>1</Version></DFSubject><DFSubject><DFSUFTCode>ON</DFSUFTCode><DFSUDescriptor>General_News</DFSUDescriptor><Version>1</Version></DFSubject></DFSubjects><DFCountries><DFCountry><DFCtryISO3166FTCode>FR</DFCtryISO3166FTCode><DFCtryISO3166Descriptor>France</DFCtryISO3166Descriptor><Version>1</Version></DFCountry></DFCountries><DFRegions><DFRegion><DFRegFTCode>XG</DFRegFTCode><DFRegDescriptor>Europe</DFRegDescriptor><Version>1</Version></DFRegion><DFRegion><DFRegFTCode>XJ</DFRegFTCode><DFRegDescriptor>Western_Europe</DFRegDescriptor><Version>1</Version></DFRegion></DFRegions><DFWPMRegions /><DFProvinces /><DFFTcomDisplayCodes /><DFFTSections /><DFWebRegions><DFWebRegion><DFWebRegFTCode>EU</DFWebRegFTCode><DFWebRegDescriptor>Europe</DFWebRegDescriptor></DFWebRegion></DFWebRegions></DataFactoryIndexing></ObjectMetadata>"

const systemAttributesWebChannel = "<props><productInfo><name>FTcom</name><issueDate>20150915</issueDate></productInfo><workFolder>/FT/WorldNews</workFolder><subFolder>UKNews</subFolder></props>"
const systemAttributesFTChannel = "<props><productInfo><name>Financial Times</name><issueDate>20150915</issueDate></productInfo><workFolder>/FT/WorldNews</workFolder><subFolder>UKNews</subFolder></props>"
const systemAttributesInvalidChannel = "<props><productInfo><name>FOOBAR</name><issueDate>20150915</issueDate></productInfo><workFolder>/FT/WorldNews</workFolder><subFolder>UKNews</subFolder></props>"
const contentWithHeadline = "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0idXRmLTgiPz4NCjwhRE9DVFlQRSBkb2MgU1lTVEVNICIvU3lzQ29uZmlnL1J1bGVzL2Z0cHNpLmR0ZCI+DQo8P0VNLWR0ZEV4dCAvU3lzQ29uZmlnL1J1bGVzL2Z0cHNpL2Z0cHNpLmR0eD8+DQo8P0VNLXRlbXBsYXRlTmFtZSAvU3lzQ29uZmlnL1RlbXBsYXRlcy9GVC9CYXNlLVN0b3J5LnhtbD8+DQo8P3htbC1mb3JtVGVtcGxhdGUgL1N5c0NvbmZpZy9UZW1wbGF0ZXMvRlQvQmFzZS1TdG9yeS54cHQ/Pg0KPD94bWwtc3R5bGVzaGVldCB0eXBlPSJ0ZXh0L2NzcyIgaHJlZj0iL1N5c0NvbmZpZy9SdWxlcy9mdHBzaS9GVC9tYWlucmVwLmNzcyI/Pg0KPGRvYyB4bWw6bGFuZz0iZW4tdWsiPjxsZWFkIGlkPSJVMzIwMDExNzk2MDc0NzhGQkUiPjxsZWFkLWhlYWRsaW5lIGlkPSJVMzIwMDExNzk2MDc0NzhhMkQiPjxuaWQtdGl0bGU+PGxuPjw/RU0tZHVtbXlUZXh0IFtJbnNlcnQgbmV3cyBpbiBkZXB0aCB0aXRsZSBoZXJlXT8+DQo8L2xuPg0KPC9uaWQtdGl0bGU+DQo8aW4tZGVwdGgtbmF2LXRpdGxlPjxsbj48P0VNLWR1bW15VGV4dCBbSW5zZXJ0IGluIGRlcHRoIG5hdiB0aXRsZSBoZXJlXT8+DQo8L2xuPg0KPC9pbi1kZXB0aC1uYXYtdGl0bGU+DQo8aGVhZGxpbmUgaWQ9IlUzMjAwMTE3OTYwNzQ3OHNtQyI+PGxuPlRoaXMgaXMgYSBoZWFkbGluZQ0KPC9sbj4NCjwvaGVhZGxpbmU+DQo8c2t5Ym94LWhlYWRsaW5lIGlkPSJVMzIwMDExNzk2MDc0Nzg1a0YiPjxsbj48P0VNLWR1bW15VGV4dCBbU2t5Ym94IGhlYWRsaW5lIGhlcmVdPz4NCjwvbG4+DQo8L3NreWJveC1oZWFkbGluZT4NCjx0cmlwbGV0LWhlYWRsaW5lPjxsbj48P0VNLWR1bW15VGV4dCBbVHJpcGxldCBoZWFkbGluZSBoZXJlXT8+DQo8L2xuPg0KPC90cmlwbGV0LWhlYWRsaW5lPg0KPHByb21vYm94LXRpdGxlPjxsbj48P0VNLWR1bW15VGV4dCBbUHJvbW9ib3ggdGl0bGUgaGVyZV0/Pg0KPC9sbj4NCjwvcHJvbW9ib3gtdGl0bGU+DQo8cHJvbW9ib3gtaGVhZGxpbmU+PGxuPjw/RU0tZHVtbXlUZXh0IFtQcm9tb2JveCBoZWFkbGluZSBoZXJlXT8+DQo8L2xuPg0KPC9wcm9tb2JveC1oZWFkbGluZT4NCjxlZGl0b3ItY2hvaWNlLWhlYWRsaW5lIGlkPSJVMzIwMDExNzk2MDc0NzhGVkUiPjxsbj48P0VNLWR1bW15VGV4dCBbU3RvcnkgcGFja2FnZSBoZWFkbGluZV0/Pg0KPC9sbj4NCjwvZWRpdG9yLWNob2ljZS1oZWFkbGluZT4NCjxuYXYtY29sbGVjdGlvbi1oZWFkbGluZT48bG4+PD9FTS1kdW1teVRleHQgW05hdiBjb2xsZWN0aW9uIGhlYWRsaW5lXT8+DQo8L2xuPg0KPC9uYXYtY29sbGVjdGlvbi1oZWFkbGluZT4NCjxpbi1kZXB0aC1uYXYtaGVhZGxpbmU+PGxuPjw/RU0tZHVtbXlUZXh0IFtJbiBkZXB0aCBuYXYgaGVhZGxpbmVdPz4NCjwvbG4+DQo8L2luLWRlcHRoLW5hdi1oZWFkbGluZT4NCjwvbGVhZC1oZWFkbGluZT4NCjx3ZWItaW5kZXgtaGVhZGxpbmUgaWQ9IlUzMjAwMTE3OTYwNzQ3OHRQRyI+PGxuPjw/RU0tZHVtbXlUZXh0IFtTaG9ydCBoZWFkbGluZV0/Pg0KPC9sbj4NCjwvd2ViLWluZGV4LWhlYWRsaW5lPg0KPHdlYi1zdGFuZC1maXJzdCBpZD0iVTMyMDAxMTc5NjA3NDc4bmVFIj48cD48P0VNLWR1bW15VGV4dCBbTG9uZyBzdGFuZGZpcnN0XT8+DQo8L3A+DQo8L3dlYi1zdGFuZC1maXJzdD4NCjx3ZWItc3ViaGVhZCBpZD0iVTMyMDAxMTc5NjA3NDc4dGpHIj48cD48P0VNLWR1bW15VGV4dCBbU2hvcnQgc3RhbmRmaXJzdF0/Pg0KPC9wPg0KPC93ZWItc3ViaGVhZD4NCjxlZGl0b3ItY2hvaWNlIGlkPSJVMzIwMDExNzk2MDc0NzhRWUUiPjw/RU0tZHVtbXlUZXh0IFtTdG9yeSBwYWNrYWdlIGh5cGVybGlua10/Pg0KPC9lZGl0b3ItY2hvaWNlPg0KPGxlYWQtdGV4dCBpZD0iVTMyMDAxMTc5NjA3NDc4T0VEIj48bGVhZC1ib2R5PjxwPjw/RU0tZHVtbXlUZXh0IFtJbnNlcnQgbGVhZCBib2R5IHRleHQgaGVyZSAtIG1pbiAxMzAgY2hhcnMsIG1heCAxNTAgY2hhcnNdPz4NCjwvcD4NCjwvbGVhZC1ib2R5Pg0KPHRyaXBsZXQtbGVhZC1ib2R5PjxwPjw/RU0tZHVtbXlUZXh0IFtJbnNlcnQgdHJpcGxldCBsZWFkIGJvZHkgYm9keSBoZXJlXT8+DQo8L3A+DQo8L3RyaXBsZXQtbGVhZC1ib2R5Pg0KPGNvbHVtbmlzdC1sZWFkLWJvZHk+PHA+PD9FTS1kdW1teVRleHQgW0luc2VydCBjb2x1bW5pc3QgbGVhZCBib2R5IGJvZHkgaGVyZV0/Pg0KPC9wPg0KPC9jb2x1bW5pc3QtbGVhZC1ib2R5Pg0KPHNob3J0LWJvZHk+PHA+PD9FTS1kdW1teVRleHQgW0luc2VydCBzaG9ydCBib2R5IGJvZHkgaGVyZV0/Pg0KPC9wPg0KPC9zaG9ydC1ib2R5Pg0KPHNreWJveC1ib2R5PjxwPjw/RU0tZHVtbXlUZXh0IFtJbnNlcnQgc2t5Ym94IGJvZHkgaGVyZV0/Pg0KPC9wPg0KPC9za3lib3gtYm9keT4NCjxwcm9tb2JveC1ib2R5PjxwPjw/RU0tZHVtbXlUZXh0IFtJbnNlcnQgcHJvbW9ib3ggYm9keSBoZXJlXT8+DQo8L3A+DQo8L3Byb21vYm94LWJvZHk+DQo8dHJpcGxldC1zaG9ydC1ib2R5PjxwPjw/RU0tZHVtbXlUZXh0IFtJbnNlcnQgdHJpcGxldCBzaG9ydCBib2R5IGhlcmVdPz4NCjwvcD4NCjwvdHJpcGxldC1zaG9ydC1ib2R5Pg0KPGVkaXRvci1jaG9pY2Utc2hvcnQtbGVhZC1ib2R5PjxwPjw/RU0tZHVtbXlUZXh0IFtJbnNlcnQgZWRpdG9ycyBjaG9pY2Ugc2hvcnQgbGVhZCBib2R5IGhlcmVdPz4NCjwvcD4NCjwvZWRpdG9yLWNob2ljZS1zaG9ydC1sZWFkLWJvZHk+DQo8bmF2LWNvbGxlY3Rpb24tc2hvcnQtbGVhZC1ib2R5PjxwPjw/RU0tZHVtbXlUZXh0IFtJbnNlcnQgbmF2IGNvbGxlY3Rpb24gc2hvcnQgbGVhZCBib2R5IGhlcmVdPz4NCjwvcD4NCjwvbmF2LWNvbGxlY3Rpb24tc2hvcnQtbGVhZC1ib2R5Pg0KPC9sZWFkLXRleHQ+DQo8bGVhZC1pbWFnZXMgaWQ9IlUzMjAwMTE3OTYwNzQ3OEFkQiI+PHdlYi1tYXN0ZXIgaWQ9IlUzMjAwMTE3OTYwNzQ3OFpSRiIvPg0KPHdlYi1za3lib3gtcGljdHVyZS8+DQo8d2ViLWFsdC1waWN0dXJlLz4NCjx3ZWItcG9wdXAtcHJldmlldyB3aWR0aD0iMTY3IiBoZWlnaHQ9Ijk2Ii8+DQo8d2ViLXBvcHVwLz4NCjwvbGVhZC1pbWFnZXM+DQo8aW50ZXJhY3RpdmUtY2hhcnQ+PD9FTS1kdW1teVRleHQgW0ludGVyYWN0aXZlLWNoYXJ0IGxpbmtdPz4NCjwvaW50ZXJhY3RpdmUtY2hhcnQ+DQo8L2xlYWQ+DQo8c3Rvcnk+PGhlYWRibG9jayBpZD0iVTMyMDAxMTc5NjA3NDc4d21IIj48aGVhZGxpbmUgaWQ9IlUzMjAwMTE3OTYwNzQ3OFgwSCI+PGxuPjw/RU0tZHVtbXlUZXh0IFtIZWFkbGluZV0/Pg0KPC9sbj4NCjwvaGVhZGxpbmU+DQo8L2hlYWRibG9jaz4NCjx0ZXh0IGlkPSJVMzIwMDExNzk2MDc0Nzh3WEQiPjxieWxpbmU+PGF1dGhvci1uYW1lPktpcmFuIFN0YWNleSwgRW5lcmd5IENvcnJlc3BvbmRlbnQ8L2F1dGhvci1uYW1lPg0KPC9ieWxpbmU+DQo8Ym9keT48cD5DT1NUIENPTVBBUklTSU9OIFdJVEggT1RIRVIgUkVHSU9OUzwvcD4NCjxwPjwvcD4NCjxwPlRvdmUgU3R1aHIgU2pvYmxvbSBzdG9vZCBpbiBmcm9udCBvZiAzMDAgb2lsIGluZHVzdHJ5IGluc2lkZXJzIGluIEFiZXJkZWVuIGxhc3Qgd2Vlaywgc2xpcHBlZCBpbnRvIHRoZSBsb2NhbCBTY290dGlzaCBkaWFsZWN0LCBhbmQgZ2F2ZSB0aGVtIG9uZSBjbGVhciBtZXNzYWdlOiDigJxLZWVwIHlvdXIgaGVpZC7igJ08L3A+DQo8cD5UaGUgTm9yd2VnaWFuIFN0YXRvaWwgZXhlY3V0aXZlIHdhcyBhZGRyZXNzaW5nIHRoZSBiaWVubmlhbCBPZmZzaG9yZSBFdXJvcGUgY29uZmVyZW5jZSwgZG9taW5hdGVkIHRoaXMgeWVhciBieSB0YWxrIG9mIHRoZSBmYWxsaW5nIG9pbCBwcmljZSBhbmQgdGhlIDxhIGhyZWY9Imh0dHA6Ly93d3cuZnQuY29tL2Ntcy9zLzAvNGYwNDQyZTAtYmNkMi0xMWU0LWE5MTctMDAxNDRmZWFiN2RlLmh0bWwjYXh6ejNsbllyTDN3YiIgdGl0bGU9Ik5vcnRoIFNlYSBvaWw6IFRoYXQgc2lua2luZyBmZWVsaW5nIC0gRlQiPmZ1dHVyZSBvZiBOb3J0aCBTZWEgcHJvZHVjdGlvbjwvYT4uIDwvcD4NCjxwPk92ZXIgdGhlIHBhc3QgMTQgbW9udGhzLCB0aGUgb2lsIHByaWNlIGhhcyBkcm9wcGVkIGZyb20gYWJvdXQgJDExNSBwZXIgYmFycmVsIHRvIGFyb3VuZCAkNTAuIE5vd2hlcmUgaGFzIHRoaXMgYmVlbiBmZWx0IG1vcmUga2Vlbmx5IHRoYW4gaW4gdGhlIE5vcnRoIFNlYSwgRk9SIERFQ0FERVMgYSBwb3dlcmhvdXNlIG9mIHRoZSBCcml0aXNoIGVjb25vbXksIGJ1dCBOT1cgYW4gYXJlYSB3aGVyZSBvaWwgZmllbGRzIGFyZSBnZXR0aW5nIG9sZGVyLCBuZXcgZGlzY292ZXJpZXMgcmFyZXIgYW5kIGNvc3RzIHN0ZWVwZXIuIE5vIHdvbmRlciBNcyBTdHVociBTam9ibG9tIGlzIHRyeWluZyB0byBzdG9wIGhlciBjb2xsZWFndWVzIGZyb20gcGFuaWNraW5nLjwvcD4NCjxwPkJ5IG5lYXJseSBhbnkgbWV0cmljLCBkb2luZyBidXNpbmVzcyBpbiB0aGUgTm9ydGggU2VhIGlzIG1vcmUgZGlmZmljdWx0IHRoYW4gYXQgYW55IG90aGVyIHBlcmlvZCBpbiBpdHMgNTAteWVhciBoaXN0b3J5LiBXaGlsZSB0aGUgb2lsIHByaWNlIGhhcyBiZWVuIGxvd2VyIOKAlCBpdCB0b3VjaGVkICQxMCBhIGJhcnJlbCBpbiAxOTg2IOKAlCBzdWNoIGEgcmFwaWQgc2x1bXAgaGFzIG5ldmVyIGJlZW4gc2VlbiB3aGVuIGNvc3RzIHdlcmUgc28gaGlnaCBhbmQgd2l0aCBvaWwgcHJvZHVjdGlvbiBhbHJlYWR5IGluIGRlY2xpbmUuPC9wPg0KPHA+SVMgUFJPRFVDVElPTiBJTiBERUNMSU5FIEJFQ0FVU0UgTk9SVEggU0VBIElTIFJVTk5JTkcgRFJZPyBCVVQgRE9FUyBUSElTIERPV05UVVJOIFJJU0sgTUFTU0lWRSBBQ0NFTEVSQVRJT04gT0YgVEhBVCBERUNMSU5FIEJFQ0FVU0UgRVhQTE9SQVRJT04vREVWRUxPUE1FTlQgSVMgQkVJTkcgQ1VSVEFJTEVEPC9wPg0KPHA+T25lIGV4ZWN1dGl2ZSBmcm9tIGFuIG9pbCBtYWpvciBzYXlzOiDigJxXZSBoYXZlIG5ldmVyIGJlZm9yZSBleHBlcmllbmNlZCBhbnl0aGluZyBsaWtlIHRoaXMgYmVmb3JlIOKAlCBldmVyeXRoaW5nIGlzIGhhcHBlbmluZyBhdCB0aGUgc2FtZSB0aW1lLuKAnTwvcD4NCjxwPkFjY29yZGluZyB0byBhbmFseXNpcyBieSBFd2FuIE11cnJheSBvZiB0aGUgY29ycG9yYXRlIGhlYWx0aCBtb25pdG9yIENvbXBhbnkgV2F0Y2gsIG91dCBvZiB0aGUgMTI2IG9pbCBleHBsb3JhdGlvbiBhbmQgcHJvZHVjdGlvbiBjb21wYW5pZXMgbGlzdGVkIGluIExvbmRvbiwgNzcgcGVyIGNlbnQgYXJlIG5vdyBsb3NzLW1ha2luZy4gQVJFIEFMTCBUSEVTRSBJTiBUSEUgTk9SVEggU0VBPz8/IFRvdGFsIGxvc3NlcyBzdGFuZCBhdCDCozYuNmJuLjwvcD4NCjxwPjxhIGhyZWY9Imh0dHA6Ly9vaWxhbmRnYXN1ay5jby51ay9lY29ub21pYy1yZXBvcnQtMjAxNS5jZm0iIHRpdGxlPSJPaWwgJmFtcDsgR2FzIEVjb25vbWljIFJlcG9ydCAyMDE1Ij5GaWd1cmVzIHJlbGVhc2VkIGxhc3Qgd2VlayBieSBPaWwgJmFtcDsgR2FzIFVLPC9hPiwgdGhlIGluZHVzdHJ5IGJvZHksIHNob3cgdGhhdCBmb3IgdGhlIGZpcnN0IHRpbWUgc2luY2UgMTk3Nywgd2hlbiBOb3J0aCBTZWEgb2lsIHdhcyBzdGlsbCB5b3VuZywgaXQgbm93IGNvc3RzIG1vcmUgdG8gcHVtcCBvaWwgb3V0IG9mIHRoZSBzZWFiZWQgdGhhbiBpdCBnZW5lcmF0ZXMgaW4gcG9zdC10YXggcmV2ZW51ZT8/PyBQUk9GSVQ/Pz8uIDwvcD4NCjxwPjxhIGhyZWY9Imh0dHA6Ly9wdWJsaWMud29vZG1hYy5jb20vcHVibGljL21lZGlhLWNlbnRyZS8xMjUyOTI1NCIgdGl0bGU9IkxvdyBvaWwgcHJpY2UgYWNjZWxlcmF0aW5nIGRlY29tbWlzc2lvbmluZyBpbiB0aGUgVUtDUyAtIFdvb2RNYWMiPmFuYWx5c2lzIHB1Ymxpc2hlZCBsYXN0IHdlZWsgYnkgV29vZCBNYWNrZW56aWU8L2E+c3VnZ2VzdHMgMTQwIGZpZWxkcyBtYXkgY2xvc2UgaW4gdGhlIG5leHQgZml2ZSB5ZWFycywgZXZlbiBpZiB0aGUgb2lsIHByaWNlIHJldHVybnMgdG8gJDg1IGEgYmFycmVsLjwvcD4NCjxwPkRpZmZlcmVudCBjb21wYW5pZXMgb3BlcmF0aW5nIGluIHRoZSBOb3J0aCBTZWEgaGF2ZSByZXNwb25kZWQgdG8gdGhpcyBpbiBkaWZmZXJlbnQgd2F5cy48L3A+DQo8cD5PaWwgbWFqb3JzIHRlbmQgdG8gaGF2ZSBhbiBleGl0IHJvdXRlLiBUaGV5IGhhdmUgY3V0IGpvYnMg4oCUIDUsNTAwIGhhdmUgYmVlbiBsb3N0IHNvIGZhciwgd2l0aCBleGVjdXRpdmVzIHdhcm5pbmcgb2YgYW5vdGhlciAxMCwwMDAgaW4gdGhlIGNvbWluZyB5ZWFycyDigJQgYW5kIGFyZSBub3cgdHJ5aW5nIHRvIFJFRFVDRSBPVEhFUiBDT1NUUyBTVUNIIEFTPz8/IGN1dCBjb3N0cy4gQnV0IGlmIGFsbCBlbHNlIGZhaWxzLCB0aGV5IGNhbiBzaW1wbHkgZXhpdCB0aGUgcmVnaW9uIGFuZCBmb2N1cyB0aGVpciBlZmZvcnRzIG9uIGNoZWFwZXIgZmllbGRzIGVsc2V3aGVyZSBpbiB0aGUgd29ybGQuPC9wPg0KPHA+RnJhbmNl4oCZcyBUb3RhbCBsYXN0IG1vbnRoIGFncmVlZCB0byBzZWxsIE5vcnRoIFNlYSA8YSBocmVmPSJodHRwOi8vd3d3LmZ0LmNvbS9jbXMvcy8wLzg5MWQ2ZTU4LTRjMTktMTFlNS1iNTU4LThhOTcyMjk3NzE4OS5odG1sI2F4enozbG5Zckwzd2IiIHRpdGxlPSJGcmVuY2ggb2lsIG1ham9yIFRvdGFsIHRvIHNlbGwgJDkwMG0gb2YgTm9ydGggU2VhIGFzc2V0cyAtIEZUIj5nYXMgdGVybWluYWxzIGFuZCBwaXBlbGluZXM8L2E+IGluIGEgJDkwMG0gZGVhbCwgd2hpbGUgRW9uLCB0aGUgR2VybWFuIHV0aWxpdHkgY29tcGFueSwgaXMgbG9va2luZyBmb3IgYnV5ZXJzIGZvciBzb21lIG9mIGl0cyBhc3NldHMgaW4gdGhlIHJlZ2lvbi48L3A+DQo8cD5TbWFsbGVyIGNvbXBhbmllcyBob3dldmVyIHRlbmQgbm90IHRvIGhhdmUgdGhhdCBvcHRpb24uIEZvciBzb21lLCB0aGUgb2lsIHByaWNlIHBsdW5nZSBtZWFucyBzcXVlZXppbmcgY29zdHMgYXMgbXVjaCBhcyBwb3NzaWJsZSBpbiBhbiBlZmZvcnQgdG8gc3RheSBhZmxvYXQgdW50aWwgaXQgcmVib3VuZHMuPC9wPg0KPHA+RE8gQ0FQRVggQ1VUIEZJUlNUIEJZIEUgTlFVRVNUPC9wPg0KPHA+VEhFTiBFWFBMQUlOIEZPQ1VTIE9OIEVYVFJBQ1RJTkcgTU9SRSBGUk9NIEZFVyBGSUVMRFMgSVQgSVMgRk9DVVNJTkcgT04sIEFORCBDVVRUSU5HIE9QRVg8L3A+DQo8cD48L3A+DQo8cD5FbnF1ZXN0LCBmb3IgZXhhbXBsZSwgaGFzIG1hZGUgYSB2aXJ0dWUgZnJvbSBnZXR0aW5nIG1vcmUgb2lsIG91dCBvZiBtYXR1cmUgZmllbGRzIHRoYW4gbWFqb3JzIGNhbi4gSXRzIFRoaXN0bGUgb2lsIGZpZWxkIHdhcyBvbmNlIG93bmVkIGJ5IEJQLCBidXQgYXMgcHJvZHVjdGlvbiBkZWNsaW5lZCwgRW5xdWVzdCBzdGVwcGVkIGluLCBhbmQgYnkgbGFzdCB5ZWFyIHRoZSBjb21wYW55IG1hbmFnZWQgdG8gZXh0cmFjdCAzbSBiYXJyZWxzIGZyb20gaXQgZm9yIHRoZSBmaXJzdCB0aW1lIHNpbmNlIDE5OTcuPC9wPg0KPHA+T25lIHdheSB0byBkbyBzbz8/Pz8gaXMgdG8gY3V0IGNvc3RzIGFnZ3Jlc3NpdmVseS4gRW5xdWVzdCBoYXMgbG9va2VkIGF0IGV2ZXJ5dGhpbmcgaW5jbHVkaW5nIHJlbG9jYXRpbmcgaXRzIGhlbGljb3B0ZXIgc2VydmljZSBpbiBhbiBlZmZvcnQgdG8gc2F2ZSBtb25leS4gQW5kIG1vcmUgaXMgbGlrZWx5IHRvIGNvbWUuPC9wPg0KPHA+QW1qYWQgQnNlc2l1LCB0aGUgY29tcGFueeKAmXMgY2hpZWYgZXhlY3V0aXZlLCBzYWlkOiDigJxUaGUgaW5kdXN0cnkgd2VudCB0aHJvdWdoIGxhc3QgeWVhcuKAmXMgY3ljbGUgaW4gYSBsaXR0bGUgYml0IG9mIGEgdGltaWQgbWFubmVyLiBUaGVyZSB3YXMgYSBmZWVsaW5nIHRoYXQgbWF5YmUgdGhpbmdzIGNhbiBjb21lIGJhY2sgbGlrZSB0aGV5IGRpZCBpbiAyMDA4IHRvIDIwMDk/Pz8sIGJ1dCB0aGF0IGhhc27igJl0IGhhcHBlbmVkLuKAnTwvcD4NCjxwPkJ1dCBmb3IgbWFueSBzbWFsbGVyIGNvbXBhbmllcywgY3V0dGluZyBjb3N0cyBoYXMgYWxzbyBpbnZvbHZlZCBjdXR0aW5nIGV4cGxvcmF0aW9uIGFuZCBkZXZlbG9wbWVudC4gT3V0IG9mIHRoZSAxMiBmaWVsZHMgaXQgb3ducyBpbiB3aGljaCBvaWwgaGFzIGJlZW4gZGlzY292ZXJlZCwgRW5xdWVzdCBpcyBjdXJyZW50bHkgZGV2ZWxvcGluZyBvbmx5IHR3by4gVGhpcyBpcyBhIHRyZW5kIG92ZXIgdGhlIGluZHVzdHJ5IGFzIGEgd2hvbGU6IGluIDIwMDggdGhlIGluZHVzdHJ5IGRyaWxsZWQgNDAgZXhwbG9yYXRpb24gd2VsbHMgaW4gdGhlIE5vcnRoIFNlYS4gTGFzdCB5ZWFyLCB0aGF0IGZpZ3VyZSB3YXMganVzdCAxNC48L3A+DQo8cD48L3A+DQo8cD5CVVQgU09NRSBDT01QQU5JRVMgQVJFIFBST0NFRURJTkcgV0lUSCBCSUcgUFJPSkVDVFM6IFNVQ0ggQVMgUFJFTUlFUiBBVCBTT0xBTjwvcD4NCjxwPklUIE1FQU5TIFNPTUUgQ09NUEFOSUVTIEFSRSBHRU5FUkFUSU5HIElOU1VGRklDSUVOVCBDQVNIIFRPIENPVkVSIENBUEVYLCBTTyBSSVNJTkcgREVCVDwvcD4NCjxwPk1hbnkgc21hbGxlciBwbGF5ZXJzIDxzcGFuIGNoYW5uZWw9IiEiPmxpa2UgRW5xdWVzdCA8L3NwYW4+aGF2ZSBhbHNvIGNvcGVkIGJ5IHRha2luZyBvbiBpbmNyZWFzaW5nIGFtb3VudHMgb2YgZGVidC48L3A+DQo8cD5QcmVtaWVyIE9pbCwgZm9yIGV4YW1wbGUsIGhhZCAkNDBtIG9mIG5ldCBjYXNoIGluIHRoZSBmaXJzdCBxdWFydGVyIG9mIDIwMDkuIEl0cyBsYXN0IHJlc3VsdHMgc2hvd2VkIHRoYXQgdGhpcyBoYXMgbm93IGJlY29tZSAkMi4xYm4gd29ydGggb2YgbmV0IGRlYnQsIGFuZCBhbmFseXN0cyBzYXkgdGhpcyB3aWxsIGtlZXAgcmlzaW5nIHVudGlsIGF0IGxlYXN0IHRoZSBlbmQgb2YgdGhlIHllYXIuIE5FVCBERUJUIFRPIEVCSVREQSBNVUxUSVBMRTwvcD4NCjxwPk92ZXIgdGhlIHNhbWUgcGVyaW9kLCB0aGUgY29tcGFueeKAmXMgb2lsIHJlc2VydmVzIGhhdmUgZmFsbGVuIHNsaWdodGx5IGZyb20gMjI4bSBiYXJyZWxzIG9mIG9pbCBlcXVpdmFsZW50IHRvIDIyM20uIFRoZSBleHRyYSBkZWJ0IGhhcyBnb25lIHRvd2FyZHMgbWFpbnRhaW5pbmcgc3RvY2tzPz8/PyByYXRoZXIgdGhhbiBncm93aW5nIHRoZW0uPC9wPg0KPHA+V2hpbGUgbWFueSBvZiB0aGVzZSBzbWFsbGVyIHBsYXllcnMgYXJlIGhpZ2hseSBpbmRlYnRlZCwgc2V2ZXJhbCBhcmUgYWxzbyB3b3JraW5nIG9uIGNvbnNpZGVyYWJsZSBuZXcgZmluZHMsIHN1Y2ggYXMgUHJlbWllcuKAmXMgU29sYW4gZmllbGQgd2VzdCBvZiBTaGV0bGFuZC4gTW9zdCBjb21wYW5pZXMgaGF2ZSBtYW5hZ2VkIHRvIHJlbmVnb3RpYXRlIHRoZWlyIGRlYnQgaW4gcmVjZW50IG1vbnRocywgbWVhbmluZyBiYW5raW5nIGNvdmVuYW50cyBhcmUgdW5saWtlbHkgdG8gYmUgYnJlYWNoZWQgaW4gdGhlIGltbWVkaWF0ZSB0ZXJtLiA8L3A+DQo8cD5BbmQgd2l0aCB0YXggYnJlYWtzIHJlY2VudGx5IG9mZmVyZWQgYnkgdGhlIFRyZWFzdXJ5LCBhbnkgaW5jcmVhc2UgaW4gcmV2ZW51ZSBpcyBsaWtlbHkgdG8gbGVhZCB0byBhIHNpbWlsYXIgaW5jcmVhc2VzIGluIHByb2ZpdHMuPC9wPg0KPHA+TWFyayBXaWxzb24sIGFuIGFuYWx5c3QgYXQgSmVmZmVyaWVzIGJhbmssIHNhaWQ6IOKAnFJpc2luZyBuZXQgZGVidCBpcyB0YWtpbmcgdmFsdWUgYXdheSBmcm9tIGVxdWl0eSBob2xkZXJzIGJ1dCBzdG9ja3MgaGF2ZSBhbiB1cHNpZGUgaWYgdGhleSBjYW4gZGVsaXZlciB3aGF0IGlzIG9uIHRoZWlyIGJhbGFuY2Ugc2hlZXQuIFdIQVQgSVMgSEUgUkVGRVJSSU5HIFRPIE9OIEJBTEFOQ0UgU0hFRVQ/4oCdPC9wPg0KPHA+PC9wPg0KPHA+QnV0IGlmIHRoZSBuZXcgZmluZHMgZGlzYXBwb2ludCwgb3IgaWYgdGhlIG9pbCBwcmljZSByZW1haW5zIHRvbyBsb3cgdG8gY292ZXIgY29zdHMgd2hlbiBvaWwgc3RhcnRzIGNvbWluZyBvdXQgb2YgdGhlIGdyb3VuZCwgdGhleSBjb3VsZCBmaW5kIHRoZW1zZWx2ZXMgZmluYW5jaWFsbHkgZGlzdHJlc3NlZC4gSU4gV0hBVCBXQVkgV0FTIEZBSVJGSUVMRCBESVNUUkVTU0VEPzwvcD4NCjxwPk9uZSBjb21wYW55IHRoYXQgaGFzIGdvbmUgdGhyb3VnaCB0aGlzIHByb2Nlc3MgaXMgRmFpcmZpZWxkLCB3aGljaCBoYXMgbWFkZSB0aGUgZGVjaXNpb24gdG8gYWJhbmRvbiBpdHMgRHVubGluIEFscGhhIHBsYXRmb3JtIGFuZCB0dXJuIGl0c2VsZiBpbnRvIGEgZGVjb21taXNzaW9uaW5nIHNwZWNpYWxpc3QgaW5zdGVhZC48L3A+DQo8cCBjaGFubmVsPSIhIj5UaGlzIGNvdWxkIGJlIGEgc291bmQgYnVzaW5lc3MgZGVjaXNpb246PC9wPg0KPHA+V2hhdGV2ZXIgaGFwcGVucywgbm9ib2R5IGludm9sdmVkIGluIE5vcnRoIFNlYSBvaWwgdGhpbmtzIGl0IHdpbGwgbG9vayB0aGUgc2FtZSBieSB0aGUgZW5kIG9mIHRoZSBkZWNhZGUuIDwvcD4NCjxwPk9uZSBzZW5pb3IgZXhlY3V0aXZlIGZyb20gYSBtYWpvciBvaWwgY29tcGFueSBzYWlkOiDigJxMb3RzIG9mIGNvbXBhbmllcyBkbyBub3QgcmVhbGlzZSBpdCB5ZXQsIGJ1dCB0aGlzIGlzIHRoZSBiZWdpbm5pbmcgb2YgdGhlIGVuZC7igJ08L3A+DQo8L2JvZHk+DQo8L3RleHQ+DQo8L3N0b3J5Pg0KPC9kb2M+DQo="
const contentWithoutHeadline = "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0idXRmLTgiPz4NCjwhRE9DVFlQRSBkb2MgU1lTVEVNICIvU3lzQ29uZmlnL1J1bGVzL2Z0cHNpLmR0ZCI+DQo8P0VNLWR0ZEV4dCAvU3lzQ29uZmlnL1J1bGVzL2Z0cHNpL2Z0cHNpLmR0eD8+DQo8P0VNLXRlbXBsYXRlTmFtZSAvU3lzQ29uZmlnL1RlbXBsYXRlcy9GVC9CYXNlLVN0b3J5LnhtbD8+DQo8P3htbC1mb3JtVGVtcGxhdGUgL1N5c0NvbmZpZy9UZW1wbGF0ZXMvRlQvQmFzZS1TdG9yeS54cHQ/Pg0KPD94bWwtc3R5bGVzaGVldCB0eXBlPSJ0ZXh0L2NzcyIgaHJlZj0iL1N5c0NvbmZpZy9SdWxlcy9mdHBzaS9GVC9tYWlucmVwLmNzcyI/Pg0KPGRvYyB4bWw6bGFuZz0iZW4tdWsiPjxsZWFkIGlkPSJVMzIwMDExNzk2MDc0NzhGQkUiPjxsZWFkLWhlYWRsaW5lIGlkPSJVMzIwMDExNzk2MDc0NzhhMkQiPjxuaWQtdGl0bGU+PGxuPjw/RU0tZHVtbXlUZXh0IFtJbnNlcnQgbmV3cyBpbiBkZXB0aCB0aXRsZSBoZXJlXT8+DQo8L2xuPg0KPC9uaWQtdGl0bGU+DQo8aW4tZGVwdGgtbmF2LXRpdGxlPjxsbj48P0VNLWR1bW15VGV4dCBbSW5zZXJ0IGluIGRlcHRoIG5hdiB0aXRsZSBoZXJlXT8+DQo8L2xuPg0KPC9pbi1kZXB0aC1uYXYtdGl0bGU+DQo8aGVhZGxpbmUgaWQ9IlUzMjAwMTE3OTYwNzQ3OHNtQyI+PGxuPg0KPC9sbj4NCjwvaGVhZGxpbmU+DQo8c2t5Ym94LWhlYWRsaW5lIGlkPSJVMzIwMDExNzk2MDc0Nzg1a0YiPjxsbj48P0VNLWR1bW15VGV4dCBbU2t5Ym94IGhlYWRsaW5lIGhlcmVdPz4NCjwvbG4+DQo8L3NreWJveC1oZWFkbGluZT4NCjx0cmlwbGV0LWhlYWRsaW5lPjxsbj48P0VNLWR1bW15VGV4dCBbVHJpcGxldCBoZWFkbGluZSBoZXJlXT8+DQo8L2xuPg0KPC90cmlwbGV0LWhlYWRsaW5lPg0KPHByb21vYm94LXRpdGxlPjxsbj48P0VNLWR1bW15VGV4dCBbUHJvbW9ib3ggdGl0bGUgaGVyZV0/Pg0KPC9sbj4NCjwvcHJvbW9ib3gtdGl0bGU+DQo8cHJvbW9ib3gtaGVhZGxpbmU+PGxuPjw/RU0tZHVtbXlUZXh0IFtQcm9tb2JveCBoZWFkbGluZSBoZXJlXT8+DQo8L2xuPg0KPC9wcm9tb2JveC1oZWFkbGluZT4NCjxlZGl0b3ItY2hvaWNlLWhlYWRsaW5lIGlkPSJVMzIwMDExNzk2MDc0NzhGVkUiPjxsbj48P0VNLWR1bW15VGV4dCBbU3RvcnkgcGFja2FnZSBoZWFkbGluZV0/Pg0KPC9sbj4NCjwvZWRpdG9yLWNob2ljZS1oZWFkbGluZT4NCjxuYXYtY29sbGVjdGlvbi1oZWFkbGluZT48bG4+PD9FTS1kdW1teVRleHQgW05hdiBjb2xsZWN0aW9uIGhlYWRsaW5lXT8+DQo8L2xuPg0KPC9uYXYtY29sbGVjdGlvbi1oZWFkbGluZT4NCjxpbi1kZXB0aC1uYXYtaGVhZGxpbmU+PGxuPjw/RU0tZHVtbXlUZXh0IFtJbiBkZXB0aCBuYXYgaGVhZGxpbmVdPz4NCjwvbG4+DQo8L2luLWRlcHRoLW5hdi1oZWFkbGluZT4NCjwvbGVhZC1oZWFkbGluZT4NCjx3ZWItaW5kZXgtaGVhZGxpbmUgaWQ9IlUzMjAwMTE3OTYwNzQ3OHRQRyI+PGxuPjw/RU0tZHVtbXlUZXh0IFtTaG9ydCBoZWFkbGluZV0/Pg0KPC9sbj4NCjwvd2ViLWluZGV4LWhlYWRsaW5lPg0KPHdlYi1zdGFuZC1maXJzdCBpZD0iVTMyMDAxMTc5NjA3NDc4bmVFIj48cD48P0VNLWR1bW15VGV4dCBbTG9uZyBzdGFuZGZpcnN0XT8+DQo8L3A+DQo8L3dlYi1zdGFuZC1maXJzdD4NCjx3ZWItc3ViaGVhZCBpZD0iVTMyMDAxMTc5NjA3NDc4dGpHIj48cD48P0VNLWR1bW15VGV4dCBbU2hvcnQgc3RhbmRmaXJzdF0/Pg0KPC9wPg0KPC93ZWItc3ViaGVhZD4NCjxlZGl0b3ItY2hvaWNlIGlkPSJVMzIwMDExNzk2MDc0NzhRWUUiPjw/RU0tZHVtbXlUZXh0IFtTdG9yeSBwYWNrYWdlIGh5cGVybGlua10/Pg0KPC9lZGl0b3ItY2hvaWNlPg0KPGxlYWQtdGV4dCBpZD0iVTMyMDAxMTc5NjA3NDc4T0VEIj48bGVhZC1ib2R5PjxwPjw/RU0tZHVtbXlUZXh0IFtJbnNlcnQgbGVhZCBib2R5IHRleHQgaGVyZSAtIG1pbiAxMzAgY2hhcnMsIG1heCAxNTAgY2hhcnNdPz4NCjwvcD4NCjwvbGVhZC1ib2R5Pg0KPHRyaXBsZXQtbGVhZC1ib2R5PjxwPjw/RU0tZHVtbXlUZXh0IFtJbnNlcnQgdHJpcGxldCBsZWFkIGJvZHkgYm9keSBoZXJlXT8+DQo8L3A+DQo8L3RyaXBsZXQtbGVhZC1ib2R5Pg0KPGNvbHVtbmlzdC1sZWFkLWJvZHk+PHA+PD9FTS1kdW1teVRleHQgW0luc2VydCBjb2x1bW5pc3QgbGVhZCBib2R5IGJvZHkgaGVyZV0/Pg0KPC9wPg0KPC9jb2x1bW5pc3QtbGVhZC1ib2R5Pg0KPHNob3J0LWJvZHk+PHA+PD9FTS1kdW1teVRleHQgW0luc2VydCBzaG9ydCBib2R5IGJvZHkgaGVyZV0/Pg0KPC9wPg0KPC9zaG9ydC1ib2R5Pg0KPHNreWJveC1ib2R5PjxwPjw/RU0tZHVtbXlUZXh0IFtJbnNlcnQgc2t5Ym94IGJvZHkgaGVyZV0/Pg0KPC9wPg0KPC9za3lib3gtYm9keT4NCjxwcm9tb2JveC1ib2R5PjxwPjw/RU0tZHVtbXlUZXh0IFtJbnNlcnQgcHJvbW9ib3ggYm9keSBoZXJlXT8+DQo8L3A+DQo8L3Byb21vYm94LWJvZHk+DQo8dHJpcGxldC1zaG9ydC1ib2R5PjxwPjw/RU0tZHVtbXlUZXh0IFtJbnNlcnQgdHJpcGxldCBzaG9ydCBib2R5IGhlcmVdPz4NCjwvcD4NCjwvdHJpcGxldC1zaG9ydC1ib2R5Pg0KPGVkaXRvci1jaG9pY2Utc2hvcnQtbGVhZC1ib2R5PjxwPjw/RU0tZHVtbXlUZXh0IFtJbnNlcnQgZWRpdG9ycyBjaG9pY2Ugc2hvcnQgbGVhZCBib2R5IGhlcmVdPz4NCjwvcD4NCjwvZWRpdG9yLWNob2ljZS1zaG9ydC1sZWFkLWJvZHk+DQo8bmF2LWNvbGxlY3Rpb24tc2hvcnQtbGVhZC1ib2R5PjxwPjw/RU0tZHVtbXlUZXh0IFtJbnNlcnQgbmF2IGNvbGxlY3Rpb24gc2hvcnQgbGVhZCBib2R5IGhlcmVdPz4NCjwvcD4NCjwvbmF2LWNvbGxlY3Rpb24tc2hvcnQtbGVhZC1ib2R5Pg0KPC9sZWFkLXRleHQ+DQo8bGVhZC1pbWFnZXMgaWQ9IlUzMjAwMTE3OTYwNzQ3OEFkQiI+PHdlYi1tYXN0ZXIgaWQ9IlUzMjAwMTE3OTYwNzQ3OFpSRiIvPg0KPHdlYi1za3lib3gtcGljdHVyZS8+DQo8d2ViLWFsdC1waWN0dXJlLz4NCjx3ZWItcG9wdXAtcHJldmlldyB3aWR0aD0iMTY3IiBoZWlnaHQ9Ijk2Ii8+DQo8d2ViLXBvcHVwLz4NCjwvbGVhZC1pbWFnZXM+DQo8aW50ZXJhY3RpdmUtY2hhcnQ+PD9FTS1kdW1teVRleHQgW0ludGVyYWN0aXZlLWNoYXJ0IGxpbmtdPz4NCjwvaW50ZXJhY3RpdmUtY2hhcnQ+DQo8L2xlYWQ+DQo8c3Rvcnk+PGhlYWRibG9jayBpZD0iVTMyMDAxMTc5NjA3NDc4d21IIj48aGVhZGxpbmUgaWQ9IlUzMjAwMTE3OTYwNzQ3OFgwSCI+PGxuPjw/RU0tZHVtbXlUZXh0IFtIZWFkbGluZV0/Pg0KPC9sbj4NCjwvaGVhZGxpbmU+DQo8L2hlYWRibG9jaz4NCjx0ZXh0IGlkPSJVMzIwMDExNzk2MDc0Nzh3WEQiPjxieWxpbmU+PGF1dGhvci1uYW1lPktpcmFuIFN0YWNleSwgRW5lcmd5IENvcnJlc3BvbmRlbnQ8L2F1dGhvci1uYW1lPg0KPC9ieWxpbmU+DQo8Ym9keT48cD5DT1NUIENPTVBBUklTSU9OIFdJVEggT1RIRVIgUkVHSU9OUzwvcD4NCjxwPjwvcD4NCjxwPlRvdmUgU3R1aHIgU2pvYmxvbSBzdG9vZCBpbiBmcm9udCBvZiAzMDAgb2lsIGluZHVzdHJ5IGluc2lkZXJzIGluIEFiZXJkZWVuIGxhc3Qgd2Vlaywgc2xpcHBlZCBpbnRvIHRoZSBsb2NhbCBTY290dGlzaCBkaWFsZWN0LCBhbmQgZ2F2ZSB0aGVtIG9uZSBjbGVhciBtZXNzYWdlOiDigJxLZWVwIHlvdXIgaGVpZC7igJ08L3A+DQo8cD5UaGUgTm9yd2VnaWFuIFN0YXRvaWwgZXhlY3V0aXZlIHdhcyBhZGRyZXNzaW5nIHRoZSBiaWVubmlhbCBPZmZzaG9yZSBFdXJvcGUgY29uZmVyZW5jZSwgZG9taW5hdGVkIHRoaXMgeWVhciBieSB0YWxrIG9mIHRoZSBmYWxsaW5nIG9pbCBwcmljZSBhbmQgdGhlIDxhIGhyZWY9Imh0dHA6Ly93d3cuZnQuY29tL2Ntcy9zLzAvNGYwNDQyZTAtYmNkMi0xMWU0LWE5MTctMDAxNDRmZWFiN2RlLmh0bWwjYXh6ejNsbllyTDN3YiIgdGl0bGU9Ik5vcnRoIFNlYSBvaWw6IFRoYXQgc2lua2luZyBmZWVsaW5nIC0gRlQiPmZ1dHVyZSBvZiBOb3J0aCBTZWEgcHJvZHVjdGlvbjwvYT4uIDwvcD4NCjxwPk92ZXIgdGhlIHBhc3QgMTQgbW9udGhzLCB0aGUgb2lsIHByaWNlIGhhcyBkcm9wcGVkIGZyb20gYWJvdXQgJDExNSBwZXIgYmFycmVsIHRvIGFyb3VuZCAkNTAuIE5vd2hlcmUgaGFzIHRoaXMgYmVlbiBmZWx0IG1vcmUga2Vlbmx5IHRoYW4gaW4gdGhlIE5vcnRoIFNlYSwgRk9SIERFQ0FERVMgYSBwb3dlcmhvdXNlIG9mIHRoZSBCcml0aXNoIGVjb25vbXksIGJ1dCBOT1cgYW4gYXJlYSB3aGVyZSBvaWwgZmllbGRzIGFyZSBnZXR0aW5nIG9sZGVyLCBuZXcgZGlzY292ZXJpZXMgcmFyZXIgYW5kIGNvc3RzIHN0ZWVwZXIuIE5vIHdvbmRlciBNcyBTdHVociBTam9ibG9tIGlzIHRyeWluZyB0byBzdG9wIGhlciBjb2xsZWFndWVzIGZyb20gcGFuaWNraW5nLjwvcD4NCjxwPkJ5IG5lYXJseSBhbnkgbWV0cmljLCBkb2luZyBidXNpbmVzcyBpbiB0aGUgTm9ydGggU2VhIGlzIG1vcmUgZGlmZmljdWx0IHRoYW4gYXQgYW55IG90aGVyIHBlcmlvZCBpbiBpdHMgNTAteWVhciBoaXN0b3J5LiBXaGlsZSB0aGUgb2lsIHByaWNlIGhhcyBiZWVuIGxvd2VyIOKAlCBpdCB0b3VjaGVkICQxMCBhIGJhcnJlbCBpbiAxOTg2IOKAlCBzdWNoIGEgcmFwaWQgc2x1bXAgaGFzIG5ldmVyIGJlZW4gc2VlbiB3aGVuIGNvc3RzIHdlcmUgc28gaGlnaCBhbmQgd2l0aCBvaWwgcHJvZHVjdGlvbiBhbHJlYWR5IGluIGRlY2xpbmUuPC9wPg0KPHA+SVMgUFJPRFVDVElPTiBJTiBERUNMSU5FIEJFQ0FVU0UgTk9SVEggU0VBIElTIFJVTk5JTkcgRFJZPyBCVVQgRE9FUyBUSElTIERPV05UVVJOIFJJU0sgTUFTU0lWRSBBQ0NFTEVSQVRJT04gT0YgVEhBVCBERUNMSU5FIEJFQ0FVU0UgRVhQTE9SQVRJT04vREVWRUxPUE1FTlQgSVMgQkVJTkcgQ1VSVEFJTEVEPC9wPg0KPHA+T25lIGV4ZWN1dGl2ZSBmcm9tIGFuIG9pbCBtYWpvciBzYXlzOiDigJxXZSBoYXZlIG5ldmVyIGJlZm9yZSBleHBlcmllbmNlZCBhbnl0aGluZyBsaWtlIHRoaXMgYmVmb3JlIOKAlCBldmVyeXRoaW5nIGlzIGhhcHBlbmluZyBhdCB0aGUgc2FtZSB0aW1lLuKAnTwvcD4NCjxwPkFjY29yZGluZyB0byBhbmFseXNpcyBieSBFd2FuIE11cnJheSBvZiB0aGUgY29ycG9yYXRlIGhlYWx0aCBtb25pdG9yIENvbXBhbnkgV2F0Y2gsIG91dCBvZiB0aGUgMTI2IG9pbCBleHBsb3JhdGlvbiBhbmQgcHJvZHVjdGlvbiBjb21wYW5pZXMgbGlzdGVkIGluIExvbmRvbiwgNzcgcGVyIGNlbnQgYXJlIG5vdyBsb3NzLW1ha2luZy4gQVJFIEFMTCBUSEVTRSBJTiBUSEUgTk9SVEggU0VBPz8/IFRvdGFsIGxvc3NlcyBzdGFuZCBhdCDCozYuNmJuLjwvcD4NCjxwPjxhIGhyZWY9Imh0dHA6Ly9vaWxhbmRnYXN1ay5jby51ay9lY29ub21pYy1yZXBvcnQtMjAxNS5jZm0iIHRpdGxlPSJPaWwgJmFtcDsgR2FzIEVjb25vbWljIFJlcG9ydCAyMDE1Ij5GaWd1cmVzIHJlbGVhc2VkIGxhc3Qgd2VlayBieSBPaWwgJmFtcDsgR2FzIFVLPC9hPiwgdGhlIGluZHVzdHJ5IGJvZHksIHNob3cgdGhhdCBmb3IgdGhlIGZpcnN0IHRpbWUgc2luY2UgMTk3Nywgd2hlbiBOb3J0aCBTZWEgb2lsIHdhcyBzdGlsbCB5b3VuZywgaXQgbm93IGNvc3RzIG1vcmUgdG8gcHVtcCBvaWwgb3V0IG9mIHRoZSBzZWFiZWQgdGhhbiBpdCBnZW5lcmF0ZXMgaW4gcG9zdC10YXggcmV2ZW51ZT8/PyBQUk9GSVQ/Pz8uIDwvcD4NCjxwPjxhIGhyZWY9Imh0dHA6Ly9wdWJsaWMud29vZG1hYy5jb20vcHVibGljL21lZGlhLWNlbnRyZS8xMjUyOTI1NCIgdGl0bGU9IkxvdyBvaWwgcHJpY2UgYWNjZWxlcmF0aW5nIGRlY29tbWlzc2lvbmluZyBpbiB0aGUgVUtDUyAtIFdvb2RNYWMiPmFuYWx5c2lzIHB1Ymxpc2hlZCBsYXN0IHdlZWsgYnkgV29vZCBNYWNrZW56aWU8L2E+c3VnZ2VzdHMgMTQwIGZpZWxkcyBtYXkgY2xvc2UgaW4gdGhlIG5leHQgZml2ZSB5ZWFycywgZXZlbiBpZiB0aGUgb2lsIHByaWNlIHJldHVybnMgdG8gJDg1IGEgYmFycmVsLjwvcD4NCjxwPkRpZmZlcmVudCBjb21wYW5pZXMgb3BlcmF0aW5nIGluIHRoZSBOb3J0aCBTZWEgaGF2ZSByZXNwb25kZWQgdG8gdGhpcyBpbiBkaWZmZXJlbnQgd2F5cy48L3A+DQo8cD5PaWwgbWFqb3JzIHRlbmQgdG8gaGF2ZSBhbiBleGl0IHJvdXRlLiBUaGV5IGhhdmUgY3V0IGpvYnMg4oCUIDUsNTAwIGhhdmUgYmVlbiBsb3N0IHNvIGZhciwgd2l0aCBleGVjdXRpdmVzIHdhcm5pbmcgb2YgYW5vdGhlciAxMCwwMDAgaW4gdGhlIGNvbWluZyB5ZWFycyDigJQgYW5kIGFyZSBub3cgdHJ5aW5nIHRvIFJFRFVDRSBPVEhFUiBDT1NUUyBTVUNIIEFTPz8/IGN1dCBjb3N0cy4gQnV0IGlmIGFsbCBlbHNlIGZhaWxzLCB0aGV5IGNhbiBzaW1wbHkgZXhpdCB0aGUgcmVnaW9uIGFuZCBmb2N1cyB0aGVpciBlZmZvcnRzIG9uIGNoZWFwZXIgZmllbGRzIGVsc2V3aGVyZSBpbiB0aGUgd29ybGQuPC9wPg0KPHA+RnJhbmNl4oCZcyBUb3RhbCBsYXN0IG1vbnRoIGFncmVlZCB0byBzZWxsIE5vcnRoIFNlYSA8YSBocmVmPSJodHRwOi8vd3d3LmZ0LmNvbS9jbXMvcy8wLzg5MWQ2ZTU4LTRjMTktMTFlNS1iNTU4LThhOTcyMjk3NzE4OS5odG1sI2F4enozbG5Zckwzd2IiIHRpdGxlPSJGcmVuY2ggb2lsIG1ham9yIFRvdGFsIHRvIHNlbGwgJDkwMG0gb2YgTm9ydGggU2VhIGFzc2V0cyAtIEZUIj5nYXMgdGVybWluYWxzIGFuZCBwaXBlbGluZXM8L2E+IGluIGEgJDkwMG0gZGVhbCwgd2hpbGUgRW9uLCB0aGUgR2VybWFuIHV0aWxpdHkgY29tcGFueSwgaXMgbG9va2luZyBmb3IgYnV5ZXJzIGZvciBzb21lIG9mIGl0cyBhc3NldHMgaW4gdGhlIHJlZ2lvbi48L3A+DQo8cD5TbWFsbGVyIGNvbXBhbmllcyBob3dldmVyIHRlbmQgbm90IHRvIGhhdmUgdGhhdCBvcHRpb24uIEZvciBzb21lLCB0aGUgb2lsIHByaWNlIHBsdW5nZSBtZWFucyBzcXVlZXppbmcgY29zdHMgYXMgbXVjaCBhcyBwb3NzaWJsZSBpbiBhbiBlZmZvcnQgdG8gc3RheSBhZmxvYXQgdW50aWwgaXQgcmVib3VuZHMuPC9wPg0KPHA+RE8gQ0FQRVggQ1VUIEZJUlNUIEJZIEUgTlFVRVNUPC9wPg0KPHA+VEhFTiBFWFBMQUlOIEZPQ1VTIE9OIEVYVFJBQ1RJTkcgTU9SRSBGUk9NIEZFVyBGSUVMRFMgSVQgSVMgRk9DVVNJTkcgT04sIEFORCBDVVRUSU5HIE9QRVg8L3A+DQo8cD48L3A+DQo8cD5FbnF1ZXN0LCBmb3IgZXhhbXBsZSwgaGFzIG1hZGUgYSB2aXJ0dWUgZnJvbSBnZXR0aW5nIG1vcmUgb2lsIG91dCBvZiBtYXR1cmUgZmllbGRzIHRoYW4gbWFqb3JzIGNhbi4gSXRzIFRoaXN0bGUgb2lsIGZpZWxkIHdhcyBvbmNlIG93bmVkIGJ5IEJQLCBidXQgYXMgcHJvZHVjdGlvbiBkZWNsaW5lZCwgRW5xdWVzdCBzdGVwcGVkIGluLCBhbmQgYnkgbGFzdCB5ZWFyIHRoZSBjb21wYW55IG1hbmFnZWQgdG8gZXh0cmFjdCAzbSBiYXJyZWxzIGZyb20gaXQgZm9yIHRoZSBmaXJzdCB0aW1lIHNpbmNlIDE5OTcuPC9wPg0KPHA+T25lIHdheSB0byBkbyBzbz8/Pz8gaXMgdG8gY3V0IGNvc3RzIGFnZ3Jlc3NpdmVseS4gRW5xdWVzdCBoYXMgbG9va2VkIGF0IGV2ZXJ5dGhpbmcgaW5jbHVkaW5nIHJlbG9jYXRpbmcgaXRzIGhlbGljb3B0ZXIgc2VydmljZSBpbiBhbiBlZmZvcnQgdG8gc2F2ZSBtb25leS4gQW5kIG1vcmUgaXMgbGlrZWx5IHRvIGNvbWUuPC9wPg0KPHA+QW1qYWQgQnNlc2l1LCB0aGUgY29tcGFueeKAmXMgY2hpZWYgZXhlY3V0aXZlLCBzYWlkOiDigJxUaGUgaW5kdXN0cnkgd2VudCB0aHJvdWdoIGxhc3QgeWVhcuKAmXMgY3ljbGUgaW4gYSBsaXR0bGUgYml0IG9mIGEgdGltaWQgbWFubmVyLiBUaGVyZSB3YXMgYSBmZWVsaW5nIHRoYXQgbWF5YmUgdGhpbmdzIGNhbiBjb21lIGJhY2sgbGlrZSB0aGV5IGRpZCBpbiAyMDA4IHRvIDIwMDk/Pz8sIGJ1dCB0aGF0IGhhc27igJl0IGhhcHBlbmVkLuKAnTwvcD4NCjxwPkJ1dCBmb3IgbWFueSBzbWFsbGVyIGNvbXBhbmllcywgY3V0dGluZyBjb3N0cyBoYXMgYWxzbyBpbnZvbHZlZCBjdXR0aW5nIGV4cGxvcmF0aW9uIGFuZCBkZXZlbG9wbWVudC4gT3V0IG9mIHRoZSAxMiBmaWVsZHMgaXQgb3ducyBpbiB3aGljaCBvaWwgaGFzIGJlZW4gZGlzY292ZXJlZCwgRW5xdWVzdCBpcyBjdXJyZW50bHkgZGV2ZWxvcGluZyBvbmx5IHR3by4gVGhpcyBpcyBhIHRyZW5kIG92ZXIgdGhlIGluZHVzdHJ5IGFzIGEgd2hvbGU6IGluIDIwMDggdGhlIGluZHVzdHJ5IGRyaWxsZWQgNDAgZXhwbG9yYXRpb24gd2VsbHMgaW4gdGhlIE5vcnRoIFNlYS4gTGFzdCB5ZWFyLCB0aGF0IGZpZ3VyZSB3YXMganVzdCAxNC48L3A+DQo8cD48L3A+DQo8cD5CVVQgU09NRSBDT01QQU5JRVMgQVJFIFBST0NFRURJTkcgV0lUSCBCSUcgUFJPSkVDVFM6IFNVQ0ggQVMgUFJFTUlFUiBBVCBTT0xBTjwvcD4NCjxwPklUIE1FQU5TIFNPTUUgQ09NUEFOSUVTIEFSRSBHRU5FUkFUSU5HIElOU1VGRklDSUVOVCBDQVNIIFRPIENPVkVSIENBUEVYLCBTTyBSSVNJTkcgREVCVDwvcD4NCjxwPk1hbnkgc21hbGxlciBwbGF5ZXJzIDxzcGFuIGNoYW5uZWw9IiEiPmxpa2UgRW5xdWVzdCA8L3NwYW4+aGF2ZSBhbHNvIGNvcGVkIGJ5IHRha2luZyBvbiBpbmNyZWFzaW5nIGFtb3VudHMgb2YgZGVidC48L3A+DQo8cD5QcmVtaWVyIE9pbCwgZm9yIGV4YW1wbGUsIGhhZCAkNDBtIG9mIG5ldCBjYXNoIGluIHRoZSBmaXJzdCBxdWFydGVyIG9mIDIwMDkuIEl0cyBsYXN0IHJlc3VsdHMgc2hvd2VkIHRoYXQgdGhpcyBoYXMgbm93IGJlY29tZSAkMi4xYm4gd29ydGggb2YgbmV0IGRlYnQsIGFuZCBhbmFseXN0cyBzYXkgdGhpcyB3aWxsIGtlZXAgcmlzaW5nIHVudGlsIGF0IGxlYXN0IHRoZSBlbmQgb2YgdGhlIHllYXIuIE5FVCBERUJUIFRPIEVCSVREQSBNVUxUSVBMRTwvcD4NCjxwPk92ZXIgdGhlIHNhbWUgcGVyaW9kLCB0aGUgY29tcGFueeKAmXMgb2lsIHJlc2VydmVzIGhhdmUgZmFsbGVuIHNsaWdodGx5IGZyb20gMjI4bSBiYXJyZWxzIG9mIG9pbCBlcXVpdmFsZW50IHRvIDIyM20uIFRoZSBleHRyYSBkZWJ0IGhhcyBnb25lIHRvd2FyZHMgbWFpbnRhaW5pbmcgc3RvY2tzPz8/PyByYXRoZXIgdGhhbiBncm93aW5nIHRoZW0uPC9wPg0KPHA+V2hpbGUgbWFueSBvZiB0aGVzZSBzbWFsbGVyIHBsYXllcnMgYXJlIGhpZ2hseSBpbmRlYnRlZCwgc2V2ZXJhbCBhcmUgYWxzbyB3b3JraW5nIG9uIGNvbnNpZGVyYWJsZSBuZXcgZmluZHMsIHN1Y2ggYXMgUHJlbWllcuKAmXMgU29sYW4gZmllbGQgd2VzdCBvZiBTaGV0bGFuZC4gTW9zdCBjb21wYW5pZXMgaGF2ZSBtYW5hZ2VkIHRvIHJlbmVnb3RpYXRlIHRoZWlyIGRlYnQgaW4gcmVjZW50IG1vbnRocywgbWVhbmluZyBiYW5raW5nIGNvdmVuYW50cyBhcmUgdW5saWtlbHkgdG8gYmUgYnJlYWNoZWQgaW4gdGhlIGltbWVkaWF0ZSB0ZXJtLiA8L3A+DQo8cD5BbmQgd2l0aCB0YXggYnJlYWtzIHJlY2VudGx5IG9mZmVyZWQgYnkgdGhlIFRyZWFzdXJ5LCBhbnkgaW5jcmVhc2UgaW4gcmV2ZW51ZSBpcyBsaWtlbHkgdG8gbGVhZCB0byBhIHNpbWlsYXIgaW5jcmVhc2VzIGluIHByb2ZpdHMuPC9wPg0KPHA+TWFyayBXaWxzb24sIGFuIGFuYWx5c3QgYXQgSmVmZmVyaWVzIGJhbmssIHNhaWQ6IOKAnFJpc2luZyBuZXQgZGVidCBpcyB0YWtpbmcgdmFsdWUgYXdheSBmcm9tIGVxdWl0eSBob2xkZXJzIGJ1dCBzdG9ja3MgaGF2ZSBhbiB1cHNpZGUgaWYgdGhleSBjYW4gZGVsaXZlciB3aGF0IGlzIG9uIHRoZWlyIGJhbGFuY2Ugc2hlZXQuIFdIQVQgSVMgSEUgUkVGRVJSSU5HIFRPIE9OIEJBTEFOQ0UgU0hFRVQ/4oCdPC9wPg0KPHA+PC9wPg0KPHA+QnV0IGlmIHRoZSBuZXcgZmluZHMgZGlzYXBwb2ludCwgb3IgaWYgdGhlIG9pbCBwcmljZSByZW1haW5zIHRvbyBsb3cgdG8gY292ZXIgY29zdHMgd2hlbiBvaWwgc3RhcnRzIGNvbWluZyBvdXQgb2YgdGhlIGdyb3VuZCwgdGhleSBjb3VsZCBmaW5kIHRoZW1zZWx2ZXMgZmluYW5jaWFsbHkgZGlzdHJlc3NlZC4gSU4gV0hBVCBXQVkgV0FTIEZBSVJGSUVMRCBESVNUUkVTU0VEPzwvcD4NCjxwPk9uZSBjb21wYW55IHRoYXQgaGFzIGdvbmUgdGhyb3VnaCB0aGlzIHByb2Nlc3MgaXMgRmFpcmZpZWxkLCB3aGljaCBoYXMgbWFkZSB0aGUgZGVjaXNpb24gdG8gYWJhbmRvbiBpdHMgRHVubGluIEFscGhhIHBsYXRmb3JtIGFuZCB0dXJuIGl0c2VsZiBpbnRvIGEgZGVjb21taXNzaW9uaW5nIHNwZWNpYWxpc3QgaW5zdGVhZC48L3A+DQo8cCBjaGFubmVsPSIhIj5UaGlzIGNvdWxkIGJlIGEgc291bmQgYnVzaW5lc3MgZGVjaXNpb246PC9wPg0KPHA+V2hhdGV2ZXIgaGFwcGVucywgbm9ib2R5IGludm9sdmVkIGluIE5vcnRoIFNlYSBvaWwgdGhpbmtzIGl0IHdpbGwgbG9vayB0aGUgc2FtZSBieSB0aGUgZW5kIG9mIHRoZSBkZWNhZGUuIDwvcD4NCjxwPk9uZSBzZW5pb3IgZXhlY3V0aXZlIGZyb20gYSBtYWpvciBvaWwgY29tcGFueSBzYWlkOiDigJxMb3RzIG9mIGNvbXBhbmllcyBkbyBub3QgcmVhbGlzZSBpdCB5ZXQsIGJ1dCB0aGlzIGlzIHRoZSBiZWdpbm5pbmcgb2YgdGhlIGVuZC7igJ08L3A+DQo8L2JvZHk+DQo8L3RleHQ+DQo8L3N0b3J5Pg0KPC9kb2M+DQo="

const ImATeapot = 418
