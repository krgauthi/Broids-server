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

					JsonElement e = obj.get("t"); // Type
					if (e.isJsonPrimitive()) {
						// Now that we know it's a primitive, we know it's safe(ish) to continue
						int type = e.getAsInt();
						if (type == FRAME_SYNC) {
							System.out.println("Sync");
						} else if (type == FRAME_DELTA) {
							System.out.println("Delta");
						}
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