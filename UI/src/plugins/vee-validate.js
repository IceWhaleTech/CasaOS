import { required, confirmed, length, email, min } from "vee-validate/dist/rules";
import { extend } from "vee-validate";

extend("required", {
    ...required,
    message: "This field is required"
});

extend("email", {
    ...email,
    message: "This field must be a valid email"
});

extend("confirmed", {
    ...confirmed,
    message: "This field confirmation does not match"
});

extend("length", {
    ...length,
    message: "This field must have 2 options"
});

extend("min", {
    ...min,
    message: "This field must have more than {length} characters"
});

