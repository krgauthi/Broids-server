import com.google.gson.*;
import com.google.gson.stream.*;
import java.net.*;
import java.io.*;

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
					if (e.isJsonPrimitive()) {
						// Now that we know it's a primitive, we know it's safe(ish) to continue
						int type = e.getAsInt();
						if (type == FRAME_SYNC) {
							System.out.println("Sync");
						} else if (type == FRAME_DELTA) {
							System.out.println("Delta");
						}
					}
					e = obj.get("gt");
					if (e.isJsonPrimitive()) {
						int type = e.getAsInt();
						System.out.println("Gametype-gt =" +type);
						}
					
					JsonArray eArray;
					eArray = obj.get("d");
					Iterator<JsonElement> dataArray = eArray.iterator();

					while(dataArray.hasNext()){
						if(dataArray.next().getAsInt().equals("t")){
							int ActionType = dataArray.next().getAsInt();
							System.out.println("ActionType-t =" + ActionType);
						}
						if(dataArray.getAsInt.equals("e")){
							entityArray = dataArray.getAsJsonArray();
							Iterator<JsonElement> entityArray = entityA.iterator();
							while(entityArray.hasNext()){
								if (entityArray.next().equals("id")) {
									String id = entityArray.next();
									System.out.println("d.e.id Id-id =" +id);
								}
								if(entityArray.next().equals("t")){
									int type = entityArray.next();
									System.out.println("d.e.type Type-t =" + type);
								}
								if(entityArray.next().equals("x")){
								   float xPos = entityArray.next();
								   System.out.println("d.e.x xPos-x =" +xPos);
								}
								if(entityArray.next().equals("y")){
								   float yPos = entityArray.next();
								   System.out.println("d.e.y yPos-y =" +yPos);
								}
								if(entityArray.next().equals("d")){
								   float dPos = entityArray.next();
								   System.out.println("d.e.d dPos-d =" +dPos);
								}
								if(entityArray.next().equals("v")){
								   float vPos = entityArray.next();
								   System.out.println("d.e.v vPos-v =" +vPos);
								}

							entityArray = entityArray.next();
							}
						}	
					dataArray = dataArray.next();	
					}
				}
			}
		} catch (UnknownHostException e) {

		} catch (IOException e) {

		} catch (Exception e) {
			// Cave Johnson, we're done here.
		}
	}
}