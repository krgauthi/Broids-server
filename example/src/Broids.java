import com.google.gson.*;
import com.google.gson.stream.*;
import java.net.*;
import java.io.*;
import java.util.*;

class Broids {
	public static final int FRAME_DELTA = 1;
	public static final int FRAME_SYNC = 2;

	public static void main(String[] arg) {
		try {
			Gson g = new Gson();
			Socket s = new Socket("localhost", 9988);

			JsonObject o = new JsonObject();
			o.addProperty("g", "broids");

			JsonWriter out = new JsonWriter(new BufferedWriter(new OutputStreamWriter(s.getOutputStream())));
			g.toJson(o, out);
			out.flush();

			JsonStreamParser parser = new JsonStreamParser(new BufferedReader(new InputStreamReader(s.getInputStream())));

			JsonElement element;
			while (parser.hasNext()) {
				element = parser.next();
				if (element.isJsonObject()) {
					// Since we know we have an object, lets do what we need to with it
					JsonObject obj = element.getAsJsonObject();

					JsonElement e;
					
					e = obj.get("t"); // Type
					int frameType = e.getAsInt();
					if (frameType == FRAME_SYNC) {
						System.out.println("Sync");
					} else if (frameType == FRAME_DELTA) {
						System.out.println("Delta");
					}
					
					e = obj.get("gt");
					int time = e.getAsInt();
					System.out.println("Gametime-gt  = " + time);
					
					JsonArray eArray;
					e = obj.get("d");
					eArray = e.getAsJsonArray();
					Iterator<JsonElement> dataArray = eArray.iterator();

					while(dataArray.hasNext()){
						e = dataArray.next();
						JsonObject inner = e.getAsJsonObject();
						
						int actionType = inner.get("t").getAsInt();
						System.out.println("ActionType-t = " + actionType);
						
						JsonObject entity = inner.get("e").getAsJsonObject();

						String id = entity.get("id").getAsString();
						System.out.println("d.e.id Id-id = " +id);

						int entityType = entity.get("t").getAsInt();
						System.out.println("d.e.t Type-t = " + entityType);

						float xPos = entity.get("x").getAsFloat();
						System.out.println("d.e.x xPos-x = " +xPos);

						float yPos = entity.get("y").getAsFloat();
						System.out.println("d.e.y yPos-y = " +yPos);

						float dPos = entity.get("d").getAsFloat();
						System.out.println("d.e.d dPos-d = " +dPos);

						float vPos = entity.get("v").getAsFloat();
						System.out.println("d.e.v vPos-v = " +vPos);

					}
				}
			}
		} catch (UnknownHostException e) {

		} catch (IOException e) {

		} catch (Exception e) {
			System.out.println(e);
			// Cave Johnson, we're done here.
		}
	}
}